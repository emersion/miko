package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.MessageType;

import java.io.DataOutputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;


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

  public static FutureOutputMessage terrainRequest() {
    return (dos) -> {
      dos.writeByte(MessageType.TERRAIN_REQUEST.getId());
      // TODO terrainrequest : quels paramètres passer ?
    };
  }

  public static FutureOutputMessage entityUpdate() {
    return (dos) -> {
      dos.writeByte(MessageType.ENTITY_UPDATE.getId());
      // TODO entityupdate : quels paramètres passer ?
    };
  }

  public static FutureOutputMessage action() {
    return (dos) -> {
      dos.writeByte(MessageType.ACTION.getId());
      // TODO action : quels paramètres passer ?
    };
  }

  public static FutureOutputMessage chatSend(String message) {
    return (dos) -> {
      dos.writeByte(MessageType.CHAT_SEND.getId());
      writeString(dos, message);
    };
  }

  private static void writeString(DataOutputStream dos, String string) throws IOException {
    int size = string.length();
    if (size >= 1 << 16)
      throw new IllegalArgumentException("The string specified is too long: >= 2^16 bytes");
    dos.writeShort(size);
    byte[] data = string.getBytes(StandardCharsets.UTF_8);
    dos.write(data);
  }
}
