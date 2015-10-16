package cr.fr.saucisseroyale.miko.network;

import java.io.IOException;

/**
 * Une erreur de parse d'un message.
 *
 */
public class MessageParseException extends IOException {

  private static final long serialVersionUID = 403825774315632241L;

  public MessageParseException() {}

  public MessageParseException(String message) {
    super(message);
  }

  public MessageParseException(String message, Throwable cause) {
    super(message, cause);
  }

  public MessageParseException(Throwable cause) {
    super(cause);
  }

}
