package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.Miko;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.EntityType;
import cr.fr.saucisseroyale.miko.protocol.EntityUpdateType;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.MessageType;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.protocol.TerrainPoint;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.io.ByteArrayOutputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


/**
 * Factory permettant de créer des {@link FutureOutputMessage} à partir de méthodes haut niveau.
 *
 */
public class OutputMessageFactory {

  private static Logger logger = LogManager.getLogger("miko.output");

  // Classe statique
  private OutputMessageFactory() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }

  public static FutureOutputMessage ping() {
    return (dos) -> {
      logger.trace("Sent ping");
      dos.writeByte(MessageType.PING.getId());
    };
  }

  public static FutureOutputMessage pong() {
    return (dos) -> {
      logger.trace("Sent pong");
      dos.writeByte(MessageType.PONG.getId());
    };
  }

  public static FutureOutputMessage exit(ExitType type) {
    return (dos) -> {
      logger.trace("Sent exit");
      dos.writeByte(MessageType.EXIT.getId());
      dos.writeByte(type.getId());
    };
  }

  public static FutureOutputMessage login(String pseudo, String password) {
    return (dos) -> {
      logger.trace("Sent login");
      dos.writeByte(MessageType.LOGIN.getId());
      writeString(dos, pseudo);
      writeString(dos, password);
    };
  }

  public static FutureOutputMessage register(String pseudo, String password) {
    return (dos) -> {
      logger.trace("Sent register");
      dos.writeByte(MessageType.REGISTER.getId());
      writeString(dos, pseudo);
      writeString(dos, password);
    };
  }

  public static FutureOutputMessage terrainRequest(List<ChunkPoint> chunks) {
    // check size against list size before further processing
    int size0 = chunks.size();
    if (size0 >= 1 << 8) {
      throw new IllegalArgumentException("The specified list is too long, max size: 255 chunks");
    }
    // make defensive copy
    List<ChunkPoint> chunksCopy = new ArrayList<>(chunks.size());
    for (ChunkPoint chunk : chunks) {
      chunksCopy.add(chunk);
    }
    int size = chunksCopy.size();
    if (size >= 1 << 8) {
      throw new IllegalArgumentException("The specified list is too long, max size: 255 chunks");
    }
    return (dos) -> {
      logger.trace("Sent terrain request");
      dos.writeByte(MessageType.TERRAIN_REQUEST.getId());
      dos.writeByte(size);
      for (ChunkPoint chunk : chunksCopy) {
        dos.writeShort(chunk.getChunkX());
        dos.writeShort(chunk.getChunkY());
      }
    };
  }

  public static FutureOutputMessage entityUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    boolean[] updateTypes = new boolean[8];
    ByteArrayOutputStream baos = new ByteArrayOutputStream();
    try (DataOutputStream buffer = new DataOutputStream(baos)) {
      for (int b = 0; b < 8; b++) {
        EntityUpdateType updateType = EntityUpdateType.getType(b);
        if (updateType == null) {
          // bit doesn't exist in the protocol or isn't implemented in client _at all_
          // ignore
          continue;
        }
        switch (updateType) {
          case POSITION:
            if (!entityDataUpdate.hasPosition()) {
              break;
            }
            updateTypes[b] = true;
            TerrainPoint point = entityDataUpdate.getPosition().toTerrainPoint();
            writeTerrainPoint(buffer, point);
            break;
          case SPEED_ANGLE:
            if (!entityDataUpdate.hasSpeedAngle()) {
              break;
            }
            updateTypes[b] = true;
            float speedAngle = entityDataUpdate.getSpeedAngle();
            buffer.writeFloat(speedAngle);
            break;
          case SPEED_NORM:
            if (!entityDataUpdate.hasSpeedNorm()) {
              break;
            }
            updateTypes[b] = true;
            float speedNorm = entityDataUpdate.getSpeedNorm();
            buffer.writeFloat(speedNorm);
            break;
          case ENTITY_TYPE:
            if (!entityDataUpdate.hasEntityType()) {
              break;
            }
            updateTypes[b] = true;
            EntityType entityType = entityDataUpdate.getEntityType();
            buffer.writeShort(entityType.getId());
            break;
          case SPRITE_TYPE:
            if (!entityDataUpdate.hasSprite()) {
              break;
            }
            updateTypes[b] = true;
            SpriteType spriteType = entityDataUpdate.getSpriteType();
            buffer.writeShort(spriteType.getId());
            break;
          case OBJECT_DATA:
            Map<ObjectAttribute, Object> attributes = entityDataUpdate.getObjectAttributes();
            if (attributes.isEmpty()) {
              break;
            }
            updateTypes[b] = true;
            int size = attributes.size();
            if (size >= 1 << 16) {
              throw new IllegalArgumentException("The attributes list is too long, max size : 65535 attributes");
            }
            buffer.writeShort(size);
            for (Map.Entry<ObjectAttribute, Object> attribute : attributes.entrySet()) {
              ObjectAttribute type = attribute.getKey();
              Object value = attribute.getValue();
              buffer.writeShort(type.getId());
              switch (type.getDataType()) {
                case VOID:
                  break;
                case ONE_SHORT:
                case ONE_ENTITY:
                  buffer.writeShort((int) value);
                  break;
                case ONE_FLOAT:
                  buffer.writeFloat((float) value);
                  break;
                case ONE_TERRAIN:
                  writeTerrainPoint(buffer, (TerrainPoint) value);
                  break;
                case PAIR_FLOAT_ENTITY:
                  @SuppressWarnings("unchecked")
                  Pair<Float, Integer> pair = (Pair<Float, Integer>) value;
                  buffer.writeFloat(pair.getFirst());
                  buffer.writeShort(pair.getSecond());
                  break;
                default:
                  throw new IllegalArgumentException("Unknown parameters set in entityDataUpdate objectattributes");
              }
            }
            break;
          default:
            throw new IllegalArgumentException("Unknown parameters set in entityDataUpdate");
        }
      }
    } catch (IOException e) {
      // impossible to get here since we write to a local byte array stream
      throw new RuntimeException(e);
    }
    int updateTypesByte = bitfieldToByte(updateTypes);
    return (dos) -> {
      logger.trace("Sent entityupdate");
      dos.writeByte(MessageType.ENTITY_UPDATE.getId());
      dos.writeShort((int) (tick % (1 << 16)));
      dos.writeShort(entityDataUpdate.getEntityId());
      dos.writeByte(updateTypesByte);
      dos.write(baos.toByteArray());
    };
  }

  public static FutureOutputMessage action(long tick, Action action) {
    return (dos) -> {
      logger.trace("Sent action");
      dos.writeByte(MessageType.ACTION.getId());
      dos.writeShort((int) (tick % (1 << 16)));
      writeAction(dos, action);
    };
  }

  public static FutureOutputMessage chatSend(String message) {
    return (dos) -> {
      logger.trace("Sent chatsend");
      dos.writeByte(MessageType.CHAT_SEND.getId());
      writeString(dos, message);
    };
  }

  public static FutureOutputMessage version() {
    return (dos) -> {
      logger.trace("Sent version");
      dos.writeByte(MessageType.VERSION.getId());
      dos.writeShort(Miko.PROTOCOL_VERSION);
    };
  }

  private static int bitfieldToByte(boolean[] values) {
    if (values.length != 8) {
      throw new IllegalArgumentException("The specified array has to be 8 booleans long");
    }
    int b = 0;
    for (int i = 0; i < 8; i++) {
      if (values[i]) {
        b |= 1 << 7 - i; // on met à 1 le bit en position i en partant de la gauche
      }
    }
    return b;
  }

  private static void writeAction(DataOutputStream dos, Action action) throws IOException {
    dos.writeShort(action.getType().getId());
    switch (action.getDataType()) {
      case VOID:
        break;
      case ONE_FLOAT:
        dos.writeFloat(action.getFloatValue());
        break;
      case ONE_SHORT:
        dos.writeShort(action.getShortValue());
        break;
      case ONE_ENTITY:
        dos.writeShort(action.getEntityIdValue());
        break;
      case ONE_TERRAIN:
        writeTerrainPoint(dos, action.getTerrainPointValue());
        break;
      case PAIR_FLOAT_ENTITY:
        Pair<Float, Integer> pair = action.getPairFloatEntityIdValue();
        dos.writeFloat(pair.getFirst());
        dos.writeShort(pair.getSecond());
        break;
      default:
        throw new IllegalArgumentException("Unknown parameters type set in action");
    }
  }

  private static void writeTerrainPoint(DataOutputStream dos, TerrainPoint terrainPoint) throws IOException {
    dos.writeShort(terrainPoint.getChunkX());
    dos.writeShort(terrainPoint.getChunkY());
    dos.writeByte(terrainPoint.getBlockX());
    dos.writeByte(terrainPoint.getBlockY());
  }

  private static void writeString(DataOutputStream dos, String string) throws IOException {
    // check length by characters length first to avoid heavy data array creation
    // (length in bytes >= characters length)
    if (string.length() >= 1 << 16) {
      throw new IllegalArgumentException("The specified string is too long, max size: 65565 bytes");
    }
    byte[] data = string.getBytes(StandardCharsets.UTF_8);
    int size = data.length;
    if (size >= 1 << 16) {
      throw new IllegalArgumentException("The specified string is too long, max size: 65565 bytes");
    }
    dos.writeShort(size);
    dos.write(data);
  }
}
