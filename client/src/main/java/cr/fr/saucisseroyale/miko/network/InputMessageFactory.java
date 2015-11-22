package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.engine.Block;
import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ActionType;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.Config;
import cr.fr.saucisseroyale.miko.protocol.DataType;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.EntityType;
import cr.fr.saucisseroyale.miko.protocol.EntityUpdateType;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.LoginResponseType;
import cr.fr.saucisseroyale.miko.protocol.MessageType;
import cr.fr.saucisseroyale.miko.protocol.MetaActionType;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.RegisterResponseType;
import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.protocol.TerrainPoint;
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
  public static FutureInputMessage parseMessage(DataInputStream dis) throws MessageParseException, IOException {
    int messageCode = dis.readUnsignedByte();
    MessageType messageType = MessageType.getType(messageCode);
    if (messageType == null) {
      throw newParseException();
    }
    int tickRemainder;
    switch (messageType) {
      case PING:
        return (handler) -> handler.ping();
      case PONG:
        return (handler) -> handler.pong();
      case EXIT:
        int exitCode = dis.readUnsignedByte();
        ExitType exitType = ExitType.getType(exitCode);
        if (exitType == null) {
          throw newParseException();
        }
        return (handler) -> handler.exit(exitType);
      case LOGIN_RESPONSE:
        int loginResponseCode = dis.readUnsignedByte();
        LoginResponseType loginResponseType = LoginResponseType.getType(loginResponseCode);
        if (loginResponseType == null) {
          throw newParseException();
        }
        if (loginResponseType == LoginResponseType.OK) {
          tickRemainder = dis.readUnsignedShort();
          return (handler) -> handler.loginSuccess(tickRemainder);
        } else {
          return (handler) -> handler.loginFail(loginResponseType);
        }
      case REGISTER_RESPONSE:
        int registerResponseCode = dis.readUnsignedByte();
        RegisterResponseType registerResponseType = RegisterResponseType.getType(registerResponseCode);
        if (registerResponseType == null) {
          throw newParseException();
        }
        return (handler) -> handler.registerResponse(registerResponseType);
      case META_ACTION:
        tickRemainder = dis.readUnsignedShort();
        int entityId = dis.readUnsignedShort();
        int metaActionCode = dis.readUnsignedByte();
        MetaActionType metaActionType = MetaActionType.getType(metaActionCode);
        if (metaActionType == null) {
          throw newParseException();
        }
        switch (metaActionType) {
          case PLAYER_JOINED:
            String pseudo = readString(dis);
            return (handler) -> handler.playerJoined(tickRemainder, entityId, pseudo);
          case PLAYER_LEFT:
            return (handler) -> handler.playerLeft(tickRemainder, entityId);
          default: // unknown
            throw newParseException();
        }
      case TERRAIN_UPDATE:
        tickRemainder = dis.readUnsignedShort();
        ChunkPoint chunkPoint = readChunkPoint(dis);
        int defaultCode = dis.readUnsignedByte();
        TerrainType defaultType = TerrainType.getType(defaultCode);
        if (defaultType == null) {
          throw newParseException();
        }
        int blocksSize = dis.readUnsignedShort();
        Block[] blocks = new Block[blocksSize];
        for (int i = 0; i < blocksSize; i++) {
          Block block = readBlock(dis);
          blocks[i] = block;
        }
        Chunk chunk = new Chunk(defaultType, Arrays.asList(blocks));
        return (handler) -> handler.chunkUpdate(tickRemainder, chunkPoint, chunk);
      case ENTITIES_UPDATE:
        tickRemainder = dis.readUnsignedShort();
        int entitiesSize = dis.readUnsignedShort();
        List<EntityDataUpdate> entitiesUpdateList = new ArrayList<>(entitiesSize);
        for (int i = 0; i < entitiesSize; i++) {
          EntityDataUpdate entityDataUpdate = readEntityDataUpdate(dis);
          entitiesUpdateList.add(entityDataUpdate);
        }
        return (handler) -> handler.entitiesUpdate(tickRemainder, entitiesUpdateList);
      case ACTIONS:
        tickRemainder = dis.readUnsignedShort();
        int actionsSize = dis.readUnsignedShort();
        List<Pair<Integer, Action>> actions = new ArrayList<>(actionsSize);
        for (int i = 0; i < actionsSize; i++) {
          int id = dis.readUnsignedShort();
          Action action = readAction(dis);
          actions.add(new Pair<>(id, action));
        }
        return (handler) -> handler.actions(tickRemainder, actions);
      case ENTITY_CREATE:
        tickRemainder = dis.readUnsignedShort();
        EntityDataUpdate entityCreateUpdate = readEntityDataUpdate(dis);
        return (handler) -> {
          handler.entityCreate(tickRemainder, entityCreateUpdate);
        };
      case ENTITY_DESTROY:
        tickRemainder = dis.readUnsignedShort();
        int entityIdDestroy = dis.readUnsignedShort();
        return (handler) -> {
          handler.entityDestroy(tickRemainder, entityIdDestroy);
        };
      case CHAT_RECEIVE:
        tickRemainder = dis.readUnsignedShort();
        int entityIdChat = dis.readUnsignedShort();
        String chatMessage = readString(dis);
        return (handler) -> handler.chatReceived(tickRemainder, entityIdChat, chatMessage);
      case CONFIG:
        Config config = readConfig(dis);
        return (handler) -> handler.config(config);
      case ENTITY_ID_CHANGE:
        int oldEntityId = dis.readUnsignedShort();
        int newEntityId = dis.readUnsignedShort();
        return (handler) -> handler.entityIdChange(oldEntityId, newEntityId);
      default: // unknown or client-only
        throw newParseException();
    }
  }

  static final FutureInputMessage networkError(Exception e) {
    return (handler) -> handler.networkError(e);
  }

  private static final Action readAction(DataInputStream dis) throws IOException {
    int actionCode = dis.readUnsignedShort();
    ActionType actionType = ActionType.getType(actionCode);
    if (actionType == null) {
      throw newParseException();
    }
    Object value = readObject(dis, actionType.getDataType());
    return new Action(actionType, value);
  }

  private static final Block readBlock(DataInputStream dis) throws IOException {
    int x = dis.readUnsignedByte();
    int y = dis.readUnsignedByte();
    int code = dis.readUnsignedByte();
    TerrainType type = TerrainType.getType(code);
    if (type == null) {
      throw newParseException();
    }
    Block block = new Block(x, y, type);
    return block;
  }

  private static final ChunkPoint readChunkPoint(DataInputStream dis) throws IOException {
    int chunkX = dis.readShort();
    int chunkY = dis.readShort();
    ChunkPoint chunkPoint = new ChunkPoint(chunkX, chunkY);
    return chunkPoint;
  }

  private static final EntityDataUpdate readEntityDataUpdate(DataInputStream dis) throws IOException {
    int entityId = dis.readUnsignedShort();
    EntityDataUpdate.Builder builder = new EntityDataUpdate.Builder(entityId);
    byte bitfield = dis.readByte();
    for (int b = 0; b < 8; b++) {
      if ((bitfield & 1 << 7 - b) == 0) {
        continue;
      }
      EntityUpdateType type = EntityUpdateType.getType(b);
      if (type == null) {
        throw newParseException();
      }
      switch (type) {
        case POSITION:
          TerrainPoint terrainPoint = readTerrainPoint(dis);
          builder.position(terrainPoint);
          break;
        case SPEED_ANGLE:
          float speedAngle = dis.readFloat();
          builder.speedAngle(speedAngle);
          break;
        case SPEED_NORM:
          float speedNorm = dis.readFloat();
          builder.speedNorm(speedNorm);
          break;
        case ENTITY_TYPE:
          int entityTypeCode = dis.readUnsignedShort();
          EntityType entityType = EntityType.getType(entityTypeCode);
          if (entityType == null) {
            throw newParseException();
          }
          builder.entityType(entityType);
          break;
        case SPRITE_TYPE:
          int spriteTypeCode = dis.readUnsignedShort();
          SpriteType spriteType = SpriteType.getType(spriteTypeCode);
          if (spriteType == null) {
            throw newParseException();
          }
          builder.spriteType(spriteType);
          break;
        case OBJECT_DATA:
          int objectDataUpdateSize = dis.readUnsignedShort();
          for (int j = 0; j < objectDataUpdateSize; j++) {
            int objectAttributeCode = dis.readUnsignedShort();
            ObjectAttribute objectAttribute = ObjectAttribute.getType(objectAttributeCode);
            if (objectAttribute == null) {
              throw newParseException();
            }
            Object objectAttributeData = readObject(dis, objectAttribute.getDataType());
            builder.objectAttribute(objectAttribute, objectAttributeData);
          }
          break;
        default:
          throw newParseException();
      }
    }
    return builder.build();
  }

  private static final Object readObject(DataInputStream dis, DataType type) throws IOException {
    switch (type) {
      case VOID:
        return null;
      case ONE_FLOAT:
        return dis.readFloat();
      case ONE_SHORT:
      case ONE_ENTITY:
        return dis.readUnsignedShort();
      case ONE_TERRAIN:
        return readTerrainPoint(dis);
      case PAIR_FLOAT_ENTITY:
        float number = dis.readFloat();
        int entityId = dis.readUnsignedShort();
        return new Pair<>(number, entityId);
      default:
        throw newParseException();
    }
  }

  private static final TerrainPoint readTerrainPoint(DataInputStream dis) throws IOException {
    int chunkX = dis.readShort();
    int chunkY = dis.readShort();
    int blockX = dis.readUnsignedByte();
    int blockY = dis.readUnsignedByte();
    TerrainPoint terrainPoint = new TerrainPoint(chunkX, chunkY, blockX, blockY);
    return terrainPoint;
  }

  private static final Config readConfig(DataInputStream dis) throws IOException {
    int maxRollbackTicks = dis.readUnsignedShort();
    float defaultPlayerSpeed = dis.readFloat();
    int playerBallCooldown = dis.readUnsignedShort();
    float defaultBallSpeed = dis.readFloat();
    int defaultBallLifespan = dis.readUnsignedShort();
    return new Config(maxRollbackTicks, defaultPlayerSpeed, playerBallCooldown, defaultBallSpeed, defaultBallLifespan);
  }

  // On utilise notre propre méthode de lecture de String au cas où le protocole change (au lieu
  // d'utiliser dis.readUTF() qui par coïncidence utilise la même méthode que nous)
  private static final String readString(DataInputStream dis) throws IOException {
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
