package cr.fr.saucisseroyale.miko.network;

import java.io.IOException;

/**
 * Une erreur de parse d'un message.
 *
 */
public class MessageParseException extends IOException {

  // this class should not be serialized so use default id
  private static final long serialVersionUID = 1L;

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
