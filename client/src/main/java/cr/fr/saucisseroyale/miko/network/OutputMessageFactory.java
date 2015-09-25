package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.FutureOutputMessage;

/**
 * Factory permettant de créer des {@link FutureOutputMessage} à partir de méthodes haut niveau.
 *
 */
public class OutputMessageFactory {

  // Static class
  private OutputMessageFactory() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }

  // TODO créer toutes les méthodes public FutureOutputMessage

  public static FutureOutputMessage ping() {
    return (dos) -> dos.writeByte(MessageType.PING.getId(false));
  }

  public static FutureOutputMessage pong() {
    return (dos) -> dos.writeByte(MessageType.PING.getId(true));
  }

}
