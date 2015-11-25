package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.Config;
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

  public void actions(int tickRemainder, List<Pair<Integer, Action>> actions);

  public void chatReceived(int tickRemainder, int entityIdChat, String chatMessage);

  public void chunksUpdate(int tickRemainder, List<Pair<ChunkPoint, Chunk>> chunks);

  public void entityIdChange(int oldEntityId, int newEntityId);

  public void entityCreate(int tickRemainder, EntityDataUpdate entityDataUpdate);

  public void entityDestroy(int tickRemainder, int entityId);

  public void entitiesUpdate(int tickRemainder, List<EntityDataUpdate> entitiesUpdateList);

  public void exit(ExitType exitType);

  public void loginFail(LoginResponseType loginResponseType);

  public void loginSuccess(int tickRemainder);

  public void networkError(Exception e);

  public void ping();

  public void playerJoined(int tickRemainder, int entityId, String pseudo);

  public void playerLeft(int tickRemainder, int entityId);

  public void pong();

  public void registerResponse(RegisterResponseType registerResponseType);

  public void config(Config config);

}
