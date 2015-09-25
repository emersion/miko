package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.network.NetworkClient;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;

/**
 * Interface entre les messages entrants et le jeu et son engine. Fournit des méthodes appelables
 * par le message entrant pour s'exécuter.
 *
 */
public class MessageHandler {

  private NetworkClient networkClient;

  public MessageHandler(NetworkClient networkClient) {
    this.networkClient = networkClient;
    // TODO il y aura évidemment d'autres paramètres
  }

  // TODO ajouter un framework pour l'exécution de l'input message

  public void pong() {
    networkClient.putMessage(OutputMessageFactory.pong());
  }

}
