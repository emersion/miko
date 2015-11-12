package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.Config;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.awt.Point;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.stream.IntStream;

/**
 * L'engine du jeu Miko, écoutant les inputs du serveur et du client, et pouvant être affiché.
 *
 */
public class Engine {

  private static final int TICK_DIVIDER = 1 << 16;

  private Config config;

  private TerrainManager terrainManager = new TerrainManager();
  private EntityManager entityManager = new EntityManager();
  private PlayerManager playerManager = new PlayerManager();
  private ChatManager chatManager = new ChatManager();
  private TickInputManager tickInputManager = new TickInputManager();

  private List<EngineMessage> messagesBuffer;

  private int playerEntityId = -1;

  private long lastTick;

  public Engine(Config config, int tickRemainder) {
    this.config = config;
    long tick = tickRemainder;
    lastTick = tick;
    messagesBuffer = new ArrayList<>();
  }

  public void processNextTick(IntStream pressedKeys, IntStream newlyPressedKeys, Point mousePosition) {
    tickInputManager.addInput(lastTick, pressedKeys, newlyPressedKeys, mousePosition);
    // redo logic until current tick
    processBufferedMessages();
    // create next tick
    updateTickAfter(lastTick);
    // dispose old ticks
    long disposeTick = lastTick - (config.getMaxRollbackTicks() + 1); // +1 to be safe
    terrainManager.disposeUntilTick(disposeTick);
    entityManager.disposeUntilTick(disposeTick);
    playerManager.disposeUntilTick(disposeTick);
    tickInputManager.disposeUntilTick(disposeTick);
  }

  private void updateTickAfter(long tick) {
    if (tick == lastTick) {
      lastTick++;
    }
    long newTick = tick + 1;

    TickInput tickInput = tickInputManager.getInput(tick);

    // integrate positions
    entityManager.getEntities(tick).forEach((entityId) -> {
      float speedNorm = entityManager.getSpeedNorm(tick, entityId);
      if (speedNorm == 0f) {
        return;
      }
      float speedAngle = entityManager.getSpeedAngle(tick, entityId);
      BlockPoint blockPoint = entityManager.getBlockPoint(tick, entityId);
      ChunkPoint chunkPoint = entityManager.getChunkPoint(tick, entityId);
      float newX = blockPoint.getBlockX() + speedNorm * cos(speedAngle);
      float newY = blockPoint.getBlockY() + speedNorm * sin(speedAngle);
      BlockPoint newBlockPoint = new BlockPoint(newX % 256f, newY % 256f);
      ChunkPoint newChunkPoint = new ChunkPoint((int) (newX / 256), (int) (newY / 256));
      entityManager.setBlockPoint(newTick, entityId, newBlockPoint);
      if (!chunkPoint.equals(newChunkPoint)) {
        entityManager.setChunkPoint(newTick, entityId, newChunkPoint);
      }
    });

  }

  private void processAction(long tick, int entityId, Action action) {
    switch (action.getType()) {
      default:
        throw new IllegalArgumentException("Unknown action type: " + action.getType());
    }
  }

  private String formatMessage(long tick, int entityId, String chatMessage) {
    if (playerManager.isPlayer(tick, entityId)) {
      return playerManager.getPlayerName(tick, entityId) + ": " + chatMessage;
    } else {
      return "Entité inconnue: " + chatMessage;
    }
  }

  private void processBufferedMessages() {
    messagesBuffer.sort(Comparator.naturalOrder());
    long previousUpdatedTick = -1L;
    for (EngineMessage engineMessage : messagesBuffer) {
      long tick = engineMessage.getTick();
      if (previousUpdatedTick != -1 && previousUpdatedTick != tick) {
        // redo logic until this tick
        for (long updateTick = previousUpdatedTick; updateTick < tick; updateTick++) {
          updateTickAfter(updateTick);
        }
        previousUpdatedTick = tick;
      }
      Object[] data = engineMessage.getData();
      switch (engineMessage.getType()) {
        case PLAYER_JOINED:
          playerManager.addPlayer(engineMessage.getTick(), (int) data[0], (String) data[1]);
          break;
        case PLAYER_LEFT:
          playerManager.removePlayer(engineMessage.getTick(), (int) data[0]);
          break;
        case CHUNK_UPDATE:
          terrainManager.addChunk(tick, (ChunkPoint) data[0], (Chunk) data[1]);
          break;
        case ACTIONS_DONE:
          @SuppressWarnings("unchecked")
          List<Pair<Integer, Action>> actionList = (List<Pair<Integer, Action>>) data[0];
          for (Pair<Integer, Action> actionPair : actionList) {
            processAction(tick, actionPair.getFirst(), actionPair.getSecond());
          }
          break;
        case ENTITY_DESTROY:
          entityManager.destroyEntity(tick, (int) data[0]);
          break;
        case ENTITY_CREATE:
          entityManager.createEntity(tick, (EntityDataUpdate) data[0]);
          break;
        case ENTITIES_UPDATE:
          @SuppressWarnings("unchecked")
          List<EntityDataUpdate> entityDataUpdateList = (List<EntityDataUpdate>) data[0];
          for (EntityDataUpdate entityDataUpdate : entityDataUpdateList) {
            entityManager.applyUpdate(tick, entityDataUpdate);
          }
          break;
        default:
          throw new IllegalArgumentException("Unsupported engine message type : " + engineMessage.getType());
      }
    }
    messagesBuffer = new ArrayList<>();
    if (previousUpdatedTick != -1) {
      // update to current tick
      for (long updateTick = previousUpdatedTick; updateTick < lastTick; updateTick++) {
        updateTickAfter(updateTick);
      }
    }
  }

  public void actions(int tickRemainder, List<Pair<Integer, Action>> actions) {
    messagesBuffer.add(EngineMessage.newActionsMessage(getTick(tickRemainder), actions));
  }

  public void chatReceived(int tickRemainder, int entityIdChat, String chatMessage) {
    String formattedMessage = formatMessage(getTick(tickRemainder), entityIdChat, chatMessage);
    chatManager.addMessage(formattedMessage);
  }

  public void chunkUpdate(int tickRemainder, ChunkPoint chunkPoint, Chunk chunk) {
    messagesBuffer.add(EngineMessage.newChunkUpdateMessage(getTick(tickRemainder), chunkPoint, chunk));
  }

  public void entityCreate(int tickRemainder, EntityDataUpdate entityDataUpdate) {
    messagesBuffer.add(EngineMessage.newEntityCreateMessage(getTick(tickRemainder), entityDataUpdate));
  }

  public void entityDestroy(int tickRemainder, int entityId) {
    messagesBuffer.add(EngineMessage.newEntityDestroyMessage(getTick(tickRemainder), entityId));
  }

  public void entitiesUpdate(int tickRemainder, List<EntityDataUpdate> entitiesUpdateList) {
    messagesBuffer.add(EngineMessage.newEntitiesUpdateMessage(getTick(tickRemainder), entitiesUpdateList));
  }

  public void playerJoined(int tickRemainder, int entityId, String pseudo) {
    messagesBuffer.add(EngineMessage.newPlayerJoinedMessage(getTick(tickRemainder), entityId, pseudo));
  }

  public void playerLeft(int tickRemainder, int entityId) {
    messagesBuffer.add(EngineMessage.newPlayerLeftMessage(getTick(tickRemainder), entityId));
  }

  /**
   * Définit définitivement l'entityId de l'entité du joueur ; doit être appelé avant toute logique
   * ou render.
   *
   * @param entityId L'entityId du joueur.
   */
  public void setPlayerEntityId(int entityId) {
    if (playerEntityId != -1) {
      throw new IllegalStateException("setPlayerEntityId is only callable once");
    }
    playerEntityId = entityId;
  }

  private long getTick(int tickRemainder) {
    long quotient = lastTick / TICK_DIVIDER;
    long difference = tickRemainder - lastTick % TICK_DIVIDER;
    long newTick;
    boolean inverted;
    if (difference < 0) {
      inverted = true;
      difference = -difference;
    } else {
      inverted = false;
    }
    if (difference >= TICK_DIVIDER / 2) {
      if (inverted) {
        newTick = (quotient + 1) * TICK_DIVIDER + tickRemainder;
      } else {
        newTick = (quotient - 1) * TICK_DIVIDER + tickRemainder;
      }
    } else {
      newTick = quotient * TICK_DIVIDER + tickRemainder;
    }
    if (newTick > lastTick) {
      lastTick = newTick;
    }
    return newTick;
  }

  // cos and sin methods here to make it easier to switch to a table-based
  // cos/sin calculation if needed
  private static float cos(float number) {
    return (float) Math.cos(number);
  }

  private static float sin(float number) {
    return (float) Math.sin(number);
  }

}
