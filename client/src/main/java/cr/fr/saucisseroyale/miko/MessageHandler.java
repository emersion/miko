package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.protocol.*;
import cr.fr.saucisseroyale.miko.util.Pair;
import cr.fr.saucisseroyale.miko.util.Pair.Int;

import java.util.List;

/**
 * Une interface entre les messages entrants et le jeu et son engine. Les méthodes portant le nom de
 * messages doivent être appelées une fois par chaque message entrant correspondant.
 */
public interface MessageHandler {
  void actions(int tickRemainder, List<Int<Action>> actions);

  void chatReceived(int tickRemainder, int entityIdChat, String chatMessage);

  void chunksUpdate(int tickRemainder, List<Pair<ChunkPoint, Chunk>> chunks);

  void entityIdChange(int oldEntityId, int newEntityId);

  void entityCreate(int tickRemainder, EntityDataUpdate entityDataUpdate);

  void entityDestroy(int tickRemainder, int entityId);

  void entitiesUpdate(int tickRemainder, List<EntityDataUpdate> entitiesUpdateList);

  void exit(ExitType exitType);

  void loginFail(LoginResponseType loginResponseType);

  void loginSuccess(int tickRemainder, long timestamp);

  void networkError(Exception e);

  void ping();

  void playerJoined(int tickRemainder, int entityId, String pseudo);

  void playerLeft(int tickRemainder, int entityId);

  void pong();

  void registerResponse(RegisterResponseType registerResponseType);

  void config(Config config);
}
