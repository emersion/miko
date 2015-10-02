package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.MessageHandler;

/**
 * Représente un message d'entrée avec des actions à effectuer.
 * <p>
 * N'appeler qu'une seule fois {@link #execute(MessageHandler)} par message.
 *
 */
public interface FutureInputMessage {

  /**
   * Exécute la procédure stockée dans le message. N'appeler qu'une seule fois cette méthode par
   * message.
   *
   * @param handler Le {@link MessageHandler} dans lequel effectuer les actions.
   */
  public void execute(MessageHandler handler);
}
