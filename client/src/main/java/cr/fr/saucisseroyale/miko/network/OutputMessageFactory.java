package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.EntityUpdateType;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.MapPoint;
import cr.fr.saucisseroyale.miko.protocol.MessageType;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;

import java.io.ByteArrayOutputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;


/**
 * Factory permettant de créer des {@link FutureOutputMessage} à partir de méthodes haut niveau.
 *
 */
public class OutputMessageFactory {

  // Classe statique
  private OutputMessageFactory() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }

  public static FutureOutputMessage ping() {
    return (dos) -> dos.writeByte(MessageType.PING.getId());
  }

  public static FutureOutputMessage pong() {
    return (dos) -> dos.writeByte(MessageType.PONG.getId());
  }

  public static FutureOutputMessage exit(ExitType type) {
    return (dos) -> {
      dos.writeByte(MessageType.EXIT.getId());
      dos.writeByte(type.getId());
    };
  }

  public static FutureOutputMessage login(String pseudo, String password) {
    return (dos) -> {
      dos.writeByte(MessageType.LOGIN.getId());
      writeString(dos, pseudo);
      writeString(dos, password);
    };
  }

  public static FutureOutputMessage register(String pseudo, String password) {
    return (dos) -> {
      dos.writeByte(MessageType.REGISTER.getId());
      writeString(dos, pseudo);
      writeString(dos, password);
    };
  }

  public static FutureOutputMessage terrainRequest(List<ChunkPoint> chunks) {
    // check size against list size before further processing
    int size0 = chunks.size();
    if (size0 >= 1 << 8)
      throw new IllegalArgumentException("The specified list is too long, max size: 255 chunks");
    // make defensive copy
    List<ChunkPoint> chunksCopy = new ArrayList<>(chunks.size());
    for (ChunkPoint chunk : chunks)
      chunksCopy.add(chunk);
    int size = chunksCopy.size();
    if (size >= 1 << 8)
      throw new IllegalArgumentException("The specified list is too long, max size: 255 chunks");
    return (dos) -> {
      dos.writeByte(MessageType.TERRAIN_REQUEST.getId());
      dos.writeByte(size);
      for (ChunkPoint chunk : chunksCopy) {
        dos.writeShort(chunk.getChunkX());
        dos.writeShort(chunk.getChunkY());
      }
    };
  }

  public static FutureOutputMessage entityUpdate(EntityDataUpdate entityDataUpdate) {
    boolean[] updateTypes = new boolean[8];
    ByteArrayOutputStream baos = new ByteArrayOutputStream();
    DataOutputStream buffer = new DataOutputStream(baos);
    for (int b = 0; b < 8; b++) {
      EntityUpdateType updateType = EntityUpdateType.getType(b);
      if (updateType == null)
        throw new IllegalArgumentException("Unknown parameters set in entityDataUpdate");
      try {
        switch (updateType) {
          case POSITION:
            if (!entityDataUpdate.hasPosition())
              break;
            updateTypes[b] = true;
            MapPoint point = entityDataUpdate.getPosition();
            writeMapPoint(buffer, point);
            break;
          case SPEED_ANGLE:
            if (!entityDataUpdate.hasSpeedAngle())
              break;
            updateTypes[b] = true;
            float speedAngle = entityDataUpdate.getSpeedAngle();
            buffer.writeFloat(speedAngle);
            break;
          case SPEED_NORM:
            if (!entityDataUpdate.hasSpeedNorm())
              break;
            updateTypes[b] = true;
            float speedNorm = entityDataUpdate.getSpeedNorm();
            buffer.writeFloat(speedNorm);
            break;
          case OBJECT_DATA:
            if (!entityDataUpdate.hasObjectAttributes())
              break;
            updateTypes[b] = true;
            Set<ObjectAttribute> attributes = entityDataUpdate.getObjectAttributes();
            for (ObjectAttribute attribute : attributes) {
              buffer.writeByte(attribute.getId());
            }
            break;
          default:
            throw new IllegalArgumentException("Unknown parameters set in entityDataUpdate");
        }
      } catch (IOException e) {
        // impossible to get here since we write to a local byte array stream
        throw new RuntimeException(e);
      }
    }
    int updateTypesByte = bitfieldToByte(updateTypes);
    return (dos) -> {
      dos.writeByte(MessageType.ENTITY_UPDATE.getId());
      dos.writeShort(entityDataUpdate.getEntityId());
      dos.writeByte(updateTypesByte);
      dos.write(baos.toByteArray());
    };
  }



  public static FutureOutputMessage action(Action action) {
    return (dos) -> {
      dos.writeByte(MessageType.ACTION.getId());
      writeAction(dos, action);
    };
  }

  public static FutureOutputMessage chatSend(String message) {
    return (dos) -> {
      dos.writeByte(MessageType.CHAT_SEND.getId());
      writeString(dos, message);
    };
  }

  private static int bitfieldToByte(boolean[] values) {
    if (values.length != 8)
      throw new IllegalArgumentException("The specified array has to be 8 booleans long");
    int b = 0;
    for (int i = 0; i < 8; i++) {
      if (values[i])
        b |= 1 << (7 - i); // on met à 1 le bit en position i en partant de la gauche
    }
    return b;
  }

  private static void writeAction(DataOutputStream dos, Action action) throws IOException {
    dos.writeByte(action.getType().getId());
    switch (action.getParameterType()) {
      case VOID:
        break;
      case FLOAT:
        dos.writeFloat(action.getFloatValue());
        break;
      case ENTITY_ID:
        dos.writeShort(action.getEntityIdValue());
        break;
      case MAP_POINT:
        writeMapPoint(dos, action.getMapPointValue());
        break;
      default:
        throw new IllegalArgumentException("Unknown parameters type set in action");
    }
  }

  private static void writeMapPoint(DataOutputStream dos, MapPoint mapPoint) throws IOException {
    dos.writeShort(mapPoint.getChunkX());
    dos.writeShort(mapPoint.getChunkY());
    dos.writeByte(mapPoint.getBlockX());
    dos.writeByte(mapPoint.getBlockY());
  }

  private static void writeString(DataOutputStream dos, String string) throws IOException {
    // check length by characters length first to avoid heavy data array creation
    // (characters length >= length in bytes)
    if (string.length() >= 1 << 16)
      throw new IllegalArgumentException("The specified string is too long, max size: 65565 bytes");
    byte[] data = string.getBytes(StandardCharsets.UTF_8);
    int size = data.length;
    if (size >= 1 << 16)
      throw new IllegalArgumentException("The specified string is too long, max size: 65565 bytes");
    dos.writeShort(size);
    dos.write(data);
  }
}
