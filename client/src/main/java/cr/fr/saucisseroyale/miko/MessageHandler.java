package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.engine.Chunk;
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
public interface MessageHandler {

  public void actions(int tick, List<Pair<Integer, Action>> actions);

  public void chatReceived(int entityIdChat, String chatMessage);

  public void chunkUpdate(int tick, ChunkPoint chunkPoint, Chunk chunk);

  public void entityCreate(int tick, int entityId, EntityDataUpdate entityDataUpdate);

  public void entityDestroy(int tick, int entityId);

  public void entitiesUpdate(int tick, List<EntityDataUpdate> entitiesUpdateList);

  public void exit(ExitType exitType);

  public void loginFail(LoginResponseType loginResponseType);

  public void loginSuccess(int tick);

  public void ping();

  public void playerJoined(int entityId, String pseudo);

  public void playerLeft(int entityId);

  public void pong();

  public void registerResponse(RegisterResponseType registerResponseType);

}
