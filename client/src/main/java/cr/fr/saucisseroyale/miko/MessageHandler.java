package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.network.NetworkClient;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.LoginResponseType;
import cr.fr.saucisseroyale.miko.protocol.RegisterResponseType;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.List;

/**
 * Une interface entre les messages entrants et le jeu et son engine. Les méthodes portant le nom de
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
  // TODO remplir les méthodes d'entrée

  public void actions(List<Pair<Integer, Action>> actions) {
    messageReceived();
  }

  public void chatReceived(int entityIdChat, String chatMessage) {
    messageReceived();

  }

  public void chunkUpdate(ChunkPoint chunkPoint, Chunk chunk) {
    messageReceived();

  }

  public void entitiesUpdate(List<EntityDataUpdate> entitiesUpdateList) {
    messageReceived();

  }

  public void exit(ExitType exitType) {
    messageReceived();

  }

  public void loginResponse(LoginResponseType loginResponseType) {
    messageReceived();

  }

  public void ping() {
    messageReceived();
    networkClient.putMessage(OutputMessageFactory.pong());
  }

  public void playerJoined(int entityId, String pseudo) {
    messageReceived();

  }

  public void playerLeft(int entityId) {
    messageReceived();

  }

  public void pong() {
    messageReceived();
  }

  public void registerResponse(RegisterResponseType registerResponseType) {
    messageReceived();

  }

  /**
   * Hook appelé lors de la réception de n'importe quel message.
   */
  private void messageReceived() {
    // Prévenir qu'on a reçu un message pour éviter le timeout
  }

}
