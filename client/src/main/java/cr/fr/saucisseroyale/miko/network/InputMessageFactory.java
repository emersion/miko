package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.engine.Block;
import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.EntityUpdateType;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.LoginResponseType;
import cr.fr.saucisseroyale.miko.protocol.MapPoint;
import cr.fr.saucisseroyale.miko.protocol.MessageType;
import cr.fr.saucisseroyale.miko.protocol.MetaActionType;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.RegisterResponseType;
import cr.fr.saucisseroyale.miko.protocol.TerrainType;

import java.io.DataInputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * Factory permettant de parse des {@link FutureInputMessage} à partir de flux.
 *
 * @see #parseMessage(DataInputStream)
 */
class InputMessageFactory {

  private static MessageParsingException parsingException = new MessageParsingException(
      "Unknown message type, aborted parsing");

  // Classe statique
  private InputMessageFactory() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }


  /**
   * Parse un message entrant d'un flux dans un {@link FutureInputMessage}
   *
   * @param dis Le flux de données entrant pour parse un message.
   * @return Un {@link FutureInputMessage} correspondant à ce qui a été parse.
   * @throws MessageParsingException S'il y a une erreur lors de la lecture (le flux est alors dans
   *         un état <b>corrompu et irrécupérable</b>).
   * @throws IOException S'il y a une erreur quelconque lors de la récupération des données.
   */
  public static FutureInputMessage parseMessage(DataInputStream dis)
      throws MessageParsingException, IOException {
    int messageCode = dis.readUnsignedByte();
    MessageType messageType = MessageType.getType(messageCode);
    if (messageType == null)
      throw parsingException;
    switch (messageType) {
      case PING:
        return (handler) -> handler.ping();
      case PONG:
        return (handler) -> handler.pong();
      case EXIT:
        int exitCode = dis.readUnsignedByte();
        ExitType exitType = ExitType.getType(exitCode);
        if (exitType == null)
          throw parsingException;
        return (handler) -> handler.exit(exitType);
      case LOGIN_RESPONSE:
        int loginResponseCode = dis.readUnsignedByte();
        LoginResponseType loginResponseType = LoginResponseType.getType(loginResponseCode);
        if (loginResponseType == null)
          throw parsingException;
        return (handler) -> handler.loginResponse(loginResponseType);
      case REGISTER_RESPONSE:
        int registerResponseCode = dis.readUnsignedByte();
        RegisterResponseType registerResponseType =
            RegisterResponseType.getType(registerResponseCode);
        if (registerResponseType == null)
          throw parsingException;
        return (handler) -> handler.registerResponse(registerResponseType);
      case META_ACTION:
        int entityId = dis.readUnsignedShort();
        int metaActionCode = dis.readUnsignedByte();
        MetaActionType metaActionType = MetaActionType.getType(metaActionCode);
        if (metaActionType == null)
          throw parsingException;
        switch (metaActionType) {
          case PLAYER_JOINED:
            String pseudo = readString(dis);
            return (handler) -> handler.playerJoined(entityId, pseudo);
          case PLAYER_LEFT:
            return (handler) -> handler.playerLeft(entityId);
          default: // unknown
            throw new MessageParsingException("Unknown message type, aborted parsing");
        }
      case TERRAIN_UPDATE:
        int chunkX = dis.readShort();
        int chunkY = dis.readShort();
        int defaultCode = dis.readUnsignedByte();
        TerrainType defaultType = TerrainType.getType(defaultCode);
        if (defaultType == null)
          throw parsingException;
        int blocksSize = dis.readUnsignedShort();
        Block[] blocks = new Block[blocksSize];
        for (int i = 0; i < blocksSize; i++) {
          int x = dis.readUnsignedByte();
          int y = dis.readUnsignedByte();
          int code = dis.readUnsignedByte();
          TerrainType type = TerrainType.getType(code);
          if (type == null)
            throw parsingException;
          Block block = new Block(x, y, type);
          blocks[i] = block;
        }
        Chunk chunk = new Chunk(defaultType, Arrays.asList(blocks));
        return (handler) -> handler.chunkUpdate(chunkX, chunkY, chunk);
      case ENTITIES_UPDATE:
        int entitiesSize = dis.readUnsignedShort();
        List<EntityDataUpdate> entitiesUpdateList = new ArrayList<>(entitiesSize);
        for (int i = 0; i < entitiesSize; i++) {
          int entityIdUpdate = dis.readUnsignedShort();
          EntityDataUpdate dataUpdate = new EntityDataUpdate(entityIdUpdate);
          byte entitiesBitfield = dis.readByte();
          for (int b = 0; b < 8; b++) {
            if ((entitiesBitfield & 1 << (7 - b)) == 0) // si le bit numéro b est 0
              continue;
            EntityUpdateType updateType = EntityUpdateType.getType(b);
            if (updateType == null)
              throw parsingException;
            switch (updateType) {
              case POSITION:
                int entityChunkX = dis.readShort();
                int entityChunkY = dis.readShort();
                int entityX = dis.readUnsignedByte();
                int entityY = dis.readUnsignedByte();
                MapPoint entityPoint = new MapPoint(entityChunkX, entityChunkY, entityX, entityY);
                dataUpdate.setPosition(entityPoint);
                break;
              case SPEED_ANGLE:
                float entitySpeedAngle = dis.readFloat();
                dataUpdate.setSpeedAngle(entitySpeedAngle);
                break;
              case SPEED_NORM:
                float entitySpeedNorm = dis.readFloat();
                dataUpdate.setSpeedNorm(entitySpeedNorm);
                break;
              case OBJECT_DATA:
                int objectDataUpdateSize = dis.readUnsignedByte();
                List<ObjectAttribute> attributes = new ArrayList<>();
                for (int j = 0; j < objectDataUpdateSize; j++) {
                  int objectDataUpdateCode = dis.readUnsignedByte();
                  ObjectAttribute objectDataUpdateType =
                      ObjectAttribute.getType(objectDataUpdateCode);
                  if (objectDataUpdateType == null)
                    throw parsingException;
                  attributes.add(objectDataUpdateType);
                }
                dataUpdate.setObjectAttributes(attributes);
                break;
              default:
                throw parsingException;
            }
          }
          entitiesUpdateList.add(dataUpdate);
        }
        return (handler) -> handler.entitiesUpdate(entitiesUpdateList);
      case ACTIONS:
        int actionsSize = dis.readUnsignedShort();
        // TODO créer classe action, def liste actions, pour chaque, parse
        return (handler) -> {
        };
      case ENTITY_CREATE:
        int entityIdCreate = dis.readUnsignedShort();
        // TODO on lira ici
        return (handler) -> {
        };
      case ENTITY_DESTROY:
        int entityIdDestroy = dis.readUnsignedShort();
        // TODO on lira ici
        return (handler) -> {
        };
      case CHAT_RECEIVE:
        int entityIdChat = dis.readUnsignedShort();
        String chatMessage = readString(dis);
        return (handler) -> handler.chatReceived(entityIdChat, chatMessage);
      case LOGIN: // client-only
      case REGISTER: // client-only
      case TERRAIN_REQUEST: // client-only
      case ENTITY_UPDATE: // client-only
      case ACTION: // client-only
      case CHAT_SEND: // client-only
      default: // unknown
        throw parsingException;
    }
  }

  // On utilise notre propre méthode de lecture de String au cas où le protocole change (au lieu
  // d'utiliser dis.readInputStream() qui par coïncidence utilise la même méthode que nous)
  private static String readString(DataInputStream dis) throws IOException {
    int length = dis.readUnsignedShort();
    byte[] data = new byte[length];
    dis.readFully(data);
    String string = new String(data, StandardCharsets.UTF_8);
    return string;
  }

}
