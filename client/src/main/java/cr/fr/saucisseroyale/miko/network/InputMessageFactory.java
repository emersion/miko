package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.FutureInputMessage;

import java.io.DataInputStream;
import java.io.IOException;

/**
 * Factory permettant de parse des {@link FutureInputMessage} à partir de flux.
 *
 * @see #parseMessage(DataInputStream)
 */
public class InputMessageFactory {

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
   *         un état <b>corrompu</b>).
   * @throws IOException S'il y a une erreur quelconque lors de la récupération des données.
   */
  public static FutureInputMessage parseMessage(DataInputStream dis)
      throws MessageParsingException, IOException {
    int rawType = byteToInt(dis.readByte());
    MessageType type = MessageType.getType(rawType);
    switch (type) {
      // case PING:
      // return (miko)->miko.sendPing();
      // TODO ajouter cas
      default:
        throw new MessageParsingException("Unknown message type, aborted parsing");
    }
  }

  private static int byteToInt(byte b) {
    return b & 0xFF;
  }

}
