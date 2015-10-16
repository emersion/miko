package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.engine.Block;
import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ActionType;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
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
import cr.fr.saucisseroyale.miko.util.Pair;

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

  // Classe statique
  private InputMessageFactory() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }


  /**
   * Parse un message entrant d'un flux dans un {@link FutureInputMessage}
   *
   * @param dis Le flux de données entrant pour parse un message.
   * @return Un {@link FutureInputMessage} correspondant à ce qui a été parse.
   * @throws MessageParseException S'il y a une erreur lors de la lecture (le flux est alors dans un
   *         état <b>corrompu et irrécupérable</b>).
   * @throws IOException S'il y a une erreur quelconque lors de la récupération des données.
   */
  public static FutureInputMessage parseMessage(DataInputStream dis) throws MessageParseException,
      IOException {
    int messageCode = dis.readUnsignedByte();
    MessageType messageType = MessageType.getType(messageCode);
    if (messageType == null)
      throw newParseException();
    switch (messageType) {
      case PING:
        return (handler) -> handler.ping();
      case PONG:
        return (handler) -> handler.pong();
      case EXIT:
        int exitCode = dis.readUnsignedByte();
        ExitType exitType = ExitType.getType(exitCode);
        if (exitType == null)
          throw newParseException();
        return (handler) -> handler.exit(exitType);
      case LOGIN_RESPONSE:
        int loginResponseCode = dis.readUnsignedByte();
        LoginResponseType loginResponseType = LoginResponseType.getType(loginResponseCode);
        if (loginResponseType == null)
          throw newParseException();
        return (handler) -> handler.loginResponse(loginResponseType);
      case REGISTER_RESPONSE:
        int registerResponseCode = dis.readUnsignedByte();
        RegisterResponseType registerResponseType =
            RegisterResponseType.getType(registerResponseCode);
        if (registerResponseType == null)
          throw newParseException();
        return (handler) -> handler.registerResponse(registerResponseType);
      case META_ACTION:
        int entityId = dis.readUnsignedShort();
        int metaActionCode = dis.readUnsignedByte();
        MetaActionType metaActionType = MetaActionType.getType(metaActionCode);
        if (metaActionType == null)
          throw newParseException();
        switch (metaActionType) {
          case PLAYER_JOINED:
            String pseudo = readString(dis);
            return (handler) -> handler.playerJoined(entityId, pseudo);
          case PLAYER_LEFT:
            return (handler) -> handler.playerLeft(entityId);
          default: // unknown
            throw newParseException();
        }
      case TERRAIN_UPDATE:
        ChunkPoint chunkPoint = readChunkPoint(dis);
        int defaultCode = dis.readUnsignedByte();
        TerrainType defaultType = TerrainType.getType(defaultCode);
        if (defaultType == null)
          throw newParseException();
        int blocksSize = dis.readUnsignedShort();
        Block[] blocks = new Block[blocksSize];
        for (int i = 0; i < blocksSize; i++) {
          Block block = readBlock(dis);
          blocks[i] = block;
        }
        Chunk chunk = new Chunk(defaultType, Arrays.asList(blocks));
        return (handler) -> handler.chunkUpdate(chunkPoint, chunk);
      case ENTITIES_UPDATE:
        int entitiesSize = dis.readUnsignedShort();
        List<EntityDataUpdate> entitiesUpdateList = new ArrayList<>(entitiesSize);
        for (int i = 0; i < entitiesSize; i++) {
          EntityDataUpdate entityDataUpdate = readEntityDataUpdate(dis);
          entitiesUpdateList.add(entityDataUpdate);
        }
        return (handler) -> handler.entitiesUpdate(entitiesUpdateList);
      case ACTIONS:
        int actionsSize = dis.readUnsignedShort();
        List<Pair<Integer, Action>> actions = new ArrayList<>(actionsSize);
        for (int i = 0; i < actionsSize; i++) {
          int id = dis.readUnsignedShort();
          Action action = readAction(dis);
          actions.add(new Pair<>(id, action));
        }
        return (handler) -> handler.actions(actions);
      case ENTITY_CREATE:
        int entityIdCreate = dis.readUnsignedShort();
        // TODO entitycreate
        return (handler) -> {
        };
      case ENTITY_DESTROY:
        int entityIdDestroy = dis.readUnsignedShort();
        // TODO entitydestroy
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
        throw newParseException();
    }
  }

  private final static Action readAction(DataInputStream dis) throws MessageParseException,
      IOException {
    int actionCode = dis.readUnsignedShort();
    ActionType actionType = ActionType.getType(actionCode);
    if (actionType == null)
      throw newParseException();
    switch (actionType.getParametersType()) {
      case VOID:
        return new Action(actionType);
      case FLOAT:
        float floatValue = dis.readFloat();
        return new Action(actionType, floatValue);
      case ENTITY_ID:
        int entityId = dis.readUnsignedShort();
        return new Action(actionType, entityId);
      case MAP_POINT:
        MapPoint mapPoint = readMapPoint(dis);
        return new Action(actionType, mapPoint);
      default:
        throw newParseException();
    }
  }

  private final static Block readBlock(DataInputStream dis) throws MessageParseException,
      IOException {
    int x = dis.readUnsignedByte();
    int y = dis.readUnsignedByte();
    int code = dis.readUnsignedByte();
    TerrainType type = TerrainType.getType(code);
    if (type == null)
      throw newParseException();
    Block block = new Block(x, y, type);
    return block;
  }

  private final static ChunkPoint readChunkPoint(DataInputStream dis) throws MessageParseException,
      IOException {
    int chunkX = dis.readShort();
    int chunkY = dis.readShort();
    ChunkPoint chunkPoint = new ChunkPoint(chunkX, chunkY);
    return chunkPoint;
  }

  private final static EntityDataUpdate readEntityDataUpdate(DataInputStream dis)
      throws MessageParseException, IOException {
    int entityId = dis.readUnsignedShort();
    EntityDataUpdate.Builder builder = new EntityDataUpdate.Builder(entityId);
    byte bitfield = dis.readByte();
    for (int b = 0; b < 8; b++) {
      if ((bitfield & 1 << (7 - b)) == 0) // si le bit numéro b est 0
        continue;
      EntityUpdateType type = EntityUpdateType.getType(b);
      if (type == null)
        throw newParseException();
      switch (type) {
        case POSITION:
          MapPoint mapPoint = readMapPoint(dis);
          builder.position(mapPoint);
          break;
        case SPEED_ANGLE:
          float speedAngle = dis.readFloat();
          builder.speedAngle(speedAngle);
          break;
        case SPEED_NORM:
          float speedNorm = dis.readFloat();
          builder.speedNorm(speedNorm);
          break;
        case OBJECT_DATA:
          int objectDataUpdateSize = dis.readUnsignedByte();
          List<ObjectAttribute> attributes = new ArrayList<>();
          for (int j = 0; j < objectDataUpdateSize; j++) {
            int objectDataUpdateCode = dis.readUnsignedByte();
            ObjectAttribute objectDataUpdateType = ObjectAttribute.getType(objectDataUpdateCode);
            if (objectDataUpdateType == null)
              throw newParseException();
            attributes.add(objectDataUpdateType);
          }
          builder.objectAttributes(attributes);
          break;
        default:
          throw newParseException();
      }
    }
    return builder.build();
  }

  private final static MapPoint readMapPoint(DataInputStream dis) throws MessageParseException,
      IOException {
    int chunkX = dis.readShort();
    int chunkY = dis.readShort();
    int blockX = dis.readUnsignedByte();
    int blockY = dis.readUnsignedByte();
    MapPoint mapPoint = new MapPoint(chunkX, chunkY, blockX, blockY);
    return mapPoint;
  }

  // On utilise notre propre méthode de lecture de String au cas où le protocole change (au lieu
  // d'utiliser dis.readInputStream() qui par coïncidence utilise la même méthode que nous)
  private final static String readString(DataInputStream dis) throws IOException {
    int length = dis.readUnsignedShort();
    byte[] data = new byte[length];
    dis.readFully(data);
    String string = new String(data, StandardCharsets.UTF_8);
    return string;
  }

  private static MessageParseException newParseException() {
    return new MessageParseException("Unknown message type, aborted parsing");
  }

}
