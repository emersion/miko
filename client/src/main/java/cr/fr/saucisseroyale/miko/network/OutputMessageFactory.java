package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.protocol.MessageType;


/**
 * Factory permettant de créer des {@link FutureOutputMessage} à partir de méthodes haut niveau.
 *
 */
public class OutputMessageFactory {

  // Classe statique
  private OutputMessageFactory() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }

  // TODO créer toutes les méthodes public FutureOutputMessage

  public static FutureOutputMessage ping() {
    return (dos) -> dos.writeByte(MessageType.PING.getId());
  }

  public static FutureOutputMessage pong() {
    return (dos) -> dos.writeByte(MessageType.PONG.getId());
  }

}
