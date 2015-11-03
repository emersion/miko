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

class DebugMessageHandler implements MessageHandler {

  private NetworkClient networkClient;

  public DebugMessageHandler(NetworkClient networkClient) {
    this.networkClient = networkClient;
  }

  public void actions(int tick, List<Pair<Integer, Action>> actions) {
    System.out.println("actions reçues, taille : " + actions.size());
    for (Pair<Integer, Action> action : actions) {
      System.out.println("\tnuméro " + action.getFirst());
      System.out.println("\tvaleur " + action.getSecond().getType());
    }
  }

  public void chatReceived(int entityIdChat, String chatMessage) {
    System.out.println("chat received " + entityIdChat + " " + chatMessage);
  }

  public void chunkUpdate(int tick, ChunkPoint chunkPoint, Chunk chunk) {
    System.out.println("chuinkupdate " + chunkPoint.getChunkX() + " " + chunkPoint.getChunkY());
  }

  public void entityCreate(int tick, int entityId, EntityDataUpdate entityDataUpdate) {
    System.out.println("entitycreate " + entityId);
  }

  public void entityDestroy(int tick, int entityId) {
    System.out.println("entitydestroy " + entityId);
  }

  public void entitiesUpdate(int tick, List<EntityDataUpdate> entitiesUpdateList) {
    System.out.println("entitiesupdate " + entitiesUpdateList.size());
  }

  public void exit(ExitType exitType) {
    System.out.println("exit " + exitType);
  }

  public void loginFail(LoginResponseType loginResponseType) {
    System.out.println("loginfailed, error: " + loginResponseType);
  }

  public void loginSuccess(int tick) {
    System.out.println("loginsuccess tick: " + tick);
  }

  public void ping() {
    networkClient.putMessage(OutputMessageFactory.pong());
    System.out.println("ping");
  }

  public void playerJoined(int entityId, String pseudo) {
    System.out.println("player joined " + entityId + " " + pseudo);
  }

  public void playerLeft(int entityId) {
    System.out.println("playerleft " + entityId);
  }

  public void pong() {
    System.out.println("pong");
  }

  public void registerResponse(RegisterResponseType registerResponseType) {
    System.out.println("register r " + registerResponseType);

  }

}
