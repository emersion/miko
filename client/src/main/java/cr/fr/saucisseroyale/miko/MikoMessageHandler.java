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

class MikoMessageHandler implements MessageHandler {

  private NetworkClient networkClient;

  public MikoMessageHandler(NetworkClient networkClient) {
    this.networkClient = networkClient;
    // TODO il y aura évidemment d'autres paramètres
  }

  // TODO ajouter entitydestory et faire mieux entitycreate
  // TODO remplir les méthodes d'entrée

  public void actions(int tick, List<Pair<Integer, Action>> actions) {
    messageReceived();
  }

  public void chatReceived(int entityIdChat, String chatMessage) {
    messageReceived();
  }

  public void chunkUpdate(int tick, ChunkPoint chunkPoint, Chunk chunk) {
    messageReceived();
  }

  public void entityCreate(int tick, int entityId, EntityDataUpdate entityDataUpdate) {
    messageReceived();
  }

  public void entityDestroy(int tick, int entityId) {
    messageReceived();
  }

  public void entitiesUpdate(int tick, List<EntityDataUpdate> entitiesUpdateList) {
    messageReceived();
  }

  public void exit(ExitType exitType) {
    messageReceived();
  }

  public void loginFail(LoginResponseType loginResponseType) {
    messageReceived();
  }

  public void loginSuccess(int tick) {
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