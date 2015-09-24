package cr.fr.saucisseroyale.miko.network;

import java.io.IOException;

/**
 * Représente une erreur ayant eu lieu lors du parsing d'un message.
 *
 */
class MessageParsingException extends IOException {

  private static final long serialVersionUID = -7164536087094995927L;

  public MessageParsingException() {}

  public MessageParsingException(String message) {
    super(message);
  }

  public MessageParsingException(String message, Throwable cause) {
    super(message, cause);
  }

  public MessageParsingException(Throwable cause) {
    super(cause);
  }

}
