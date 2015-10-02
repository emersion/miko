package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.network.NetworkClient;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;

/**
 * Interface entre les messages entrants et le jeu et son engine. Les méthodes portant le nom de
 * messages doivent être appelées une fois par chaque message entrant correspondant.
 *
 */
public class MessageHandler {

  private NetworkClient networkClient;

  public MessageHandler(NetworkClient networkClient) {
    this.networkClient = networkClient;
    // TODO il y aura évidemment d'autres paramètres
  }

  // TODO ajouter toutes les méthodes d'entrée

  public void ping() {
    messageReceived();
    networkClient.putMessage(OutputMessageFactory.pong());
  }

  public void pong() {
    messageReceived();
  }

  /**
   * Hook appelé lors de la réception de n'importe quel message.
   */
  private void messageReceived() {
    // Prévenir qu'on a reçu un message pour éviter le timeout
  }

}
