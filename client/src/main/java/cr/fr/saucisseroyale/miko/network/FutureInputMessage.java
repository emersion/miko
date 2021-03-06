package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.MessageHandler;

/**
 * Un message d'entrée avec des actions à effectuer.
 * <p>
 * N'appeler qu'une seule fois {@link #execute(MessageHandler)} par message.
 */
@FunctionalInterface
public interface FutureInputMessage {
  /**
   * Exécute la procédure stockée dans le message. N'appeler qu'une seule fois cette méthode par
   * message.
   *
   * @param handler Le {@link MessageHandler} dans lequel effectuer les actions.
   */
  void execute(MessageHandler handler);
}
