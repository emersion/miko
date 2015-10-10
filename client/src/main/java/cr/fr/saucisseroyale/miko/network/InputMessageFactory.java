package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.engine.Block;
import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.LoginResponseType;
import cr.fr.saucisseroyale.miko.protocol.MessageType;
import cr.fr.saucisseroyale.miko.protocol.MetaActionType;
import cr.fr.saucisseroyale.miko.protocol.RegisterResponseType;
import cr.fr.saucisseroyale.miko.protocol.TerrainType;

import java.io.ByteArrayInputStream;
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
   * @throws MessageParsingException S'il y a une erreur lors de la lecture (le flux est alors dans
   *         un état <b>corrompu et irrécupérable</b>).
   * @throws IOException S'il y a une erreur quelconque lors de la récupération des données.
   */
  public static FutureInputMessage parseMessage(DataInputStream dis)
      throws MessageParsingException, IOException {
    int messageCode = dis.readUnsignedByte();
    MessageType messageType = MessageType.getType(messageCode);
    switch (messageType) {
      case PING:
        return (handler) -> handler.ping();
      case PONG:
        return (handler) -> handler.pong();
      case EXIT:
        int exitCode = dis.readUnsignedByte();
        ExitType exitType = ExitType.getType(exitCode);
        return (handler) -> handler.exit(exitType);
      case LOGIN_RESPONSE:
        int loginResponseCode = dis.readUnsignedByte();
        LoginResponseType loginResponseType = LoginResponseType.getType(loginResponseCode);
        return (handler) -> handler.loginResponse(loginResponseType);
      case REGISTER_RESPONSE:
        int registerResponseCode = dis.readUnsignedByte();
        RegisterResponseType registerResponseType =
            RegisterResponseType.getType(registerResponseCode);
        return (handler) -> handler.registerResponse(registerResponseType);
      case META_ACTION:
        int entityId = dis.readUnsignedShort();
        int metaActionCode = dis.readUnsignedByte();
        MetaActionType metaActionType = MetaActionType.getType(metaActionCode);
        switch (metaActionType) {
          case PLAYER_JOINED:
            String pseudo = readString(dis);
            return (handler) -> handler.playerJoined(entityId, pseudo);
          case PLAYER_LEFT:
            return (handler) -> handler.playerLeft(entityId);
          default:
            throw new MessageParsingException("Unknown message type, aborted parsing");
        }
      case TERRAIN_UPDATE:
        int chunkX = dis.readShort();
        int chunkY = dis.readShort();
        int defaultCode = dis.readUnsignedByte();
        TerrainType defaultType = TerrainType.getType(defaultCode);
        int blocksSize = dis.readUnsignedShort();
        Block[] blocks = new Block[blocksSize];
        for (int i = 0; i < blocksSize; i++) {
          int x = dis.readUnsignedByte();
          int y = dis.readUnsignedByte();
          int code = dis.readUnsignedByte();
          TerrainType type = TerrainType.getType(code);
          Block block = new Block(x, y, type);
          blocks[i] = block;
        }
        Chunk chunk = new Chunk(defaultType, Arrays.asList(blocks));
        return (handler) -> handler.chunkUpdate(chunkX, chunkY, chunk);
      case ENTITIES_UPDATE:
        int entitiesSize = dis.readUnsignedShort();
        List<FutureInputMessage> lambdas = new ArrayList<>();
        for (int i = 0; i < entitiesSize; i++) {
          int entityIdUpdate = dis.readUnsignedShort();
          byte entitiesBitfield = dis.readByte();
          if ((entitiesBitfield & 1 << 7) != 0) { // position update
            int entityChunkX = dis.readShort();
            int entityChunkY = dis.readShort();
            int entityX = dis.readUnsignedByte();
            int entityY = dis.readUnsignedByte();
            lambdas.add((handler) -> handler.positionUpdate(entityIdUpdate, entityChunkX,
                entityChunkY, entityX, entityY));
          }
          if ((entitiesBitfield & 1 << 6) != 0) { // speedangle update
            float entitySpeedAngle = dis.readFloat();
            lambdas.add((handler) -> handler.speedAngleUpdate(entityIdUpdate, entitySpeedAngle));
          }
          if ((entitiesBitfield & 1 << 5) != 0) { // speednorm update
            float entitySpeedNorm = dis.readFloat();
            lambdas.add((handler) -> handler.speedNormUpdate(entityIdUpdate, entitySpeedNorm));
          }
          if ((entitiesBitfield & 1 << 0) != 0) { // object update
            byte objectBitfield = dis.readByte();
            boolean[] objectBooleanField = new boolean[8];
            for (int j = 0; j < 8; j++) {
              objectBooleanField[j] = (objectBitfield & 1 << (7 - j)) != 0;
            }
            int objectSize = dis.readUnsignedShort();
            // Il vaudrait mieux le lire ici, mais difficile
            // TODO: faire lire données à entités
            byte[] objectData = new byte[objectSize];
            dis.readFully(objectData);
            DataInputStream objectInputStream =
                new DataInputStream(new ByteArrayInputStream(objectData));
            lambdas.add((handler) -> handler.objectUpdate(entityIdUpdate, objectBooleanField,
                objectSize, objectInputStream));
          }
        }
        return (handler) -> {
          for (FutureInputMessage lambda : lambdas) {
            lambda.execute(handler);
          }
        };
      case ACTIONS:
        // TODO
      case LOGIN:
        // TODO ajouter cas
      default:
        throw new MessageParsingException("Unknown message type, aborted parsing");
    }
  }

  private static String readString(DataInputStream dis) throws IOException {
    int length = dis.readUnsignedShort();
    byte[] data = new byte[length];
    dis.readFully(data);
    String string = new String(data, StandardCharsets.UTF_8);
    return string;
  }

}
