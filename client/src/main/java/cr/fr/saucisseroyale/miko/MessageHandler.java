package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.network.NetworkClient;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.LoginResponseType;
import cr.fr.saucisseroyale.miko.protocol.RegisterResponseType;

import java.io.DataInputStream;
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

  public void chatReceived(int entityIdChat, String chatMessage) {

  }

  public void chunkUpdate(int chunkX, int chunkY, Chunk chunk) {

  }

  public void exit(ExitType exitType) {

  }

  public void loginResponse(LoginResponseType loginResponseType) {

  }

  public void objectUpdate(int entityIdUpdate, boolean[] objectBooleanField, int objectSize,
      DataInputStream objectInputStream) {

  }

  public void ping() {
    messageReceived();
    networkClient.putMessage(OutputMessageFactory.pong());
  }

  public void playerJoined(int entityId, String pseudo) {

  }

  public void playerLeft(int entityId) {

  }

  public void pong() {
    messageReceived();
  }

  public void positionUpdate(int entityIdUpdate, int entityChunkX, int entityChunkY, int entityX,
      int entityY) {

  }

  public void registerResponse(RegisterResponseType registerResponseType) {

  }

  public void speedAngleUpdate(int entityIdUpdate, float entitySpeedAngle) {

  }

  public void speedNormUpdate(int entityIdUpdate, float entitySpeedNorm) {

  }

  /**
   * Hook appelé lors de la réception de n'importe quel message.
   */
  private void messageReceived() {
    // Prévenir qu'on a reçu un message pour éviter le timeout
  }

  public void entitiesUpdate(List<EntityDataUpdate> entitiesUpdateList) {

  }

}
