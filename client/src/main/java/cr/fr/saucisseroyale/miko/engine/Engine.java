package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.network.FutureOutputMessage;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;
import cr.fr.saucisseroyale.miko.protocol.*;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate.Builder;
import cr.fr.saucisseroyale.miko.util.MikoMath;
import cr.fr.saucisseroyale.miko.util.Or;
import cr.fr.saucisseroyale.miko.util.Pair;
import cr.fr.saucisseroyale.miko.util.Pair.FloatFloat;
import cr.fr.saucisseroyale.miko.util.Pair.Int;
import fr.delthas.uitest.Drawer;
import fr.delthas.uitest.Font;
import fr.delthas.uitest.Ui;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.awt.*;
import java.io.IOException;
import java.util.*;
import java.util.List;
import java.util.function.Consumer;
import java.util.function.IntConsumer;

/**
 * L'engine du jeu Miko, écoutant les inputs du serveur et du client, et pouvant être affiché.
 */
public class Engine {
  private static final int TICK_DIVIDER = 1 << 16;
  // sets how much the screen moves in response to mouse moves
  // MUST be positive
  private static final float MOUSE_SCREEN_MOVING = 0.3f;
  private static final int MAX_INTERPOLATION_DISTANCE = 100;
  private static Logger logger = LogManager.getLogger("miko.engine");
  private final Consumer<FutureOutputMessage> messageOutput;
  private Config config;
  private TerrainManager terrainManager;
  private TerrainImageManager terrainImageManager;
  private EntityManager entityManager = new EntityManager();
  private PlayerManager playerManager = new PlayerManager();
  private ChatManager chatManager = new ChatManager();
  private TickInputManager tickInputManager = new TickInputManager();
  private SpriteManager spriteManager;
  private List<EngineMessage> messagesBuffer = new LinkedList<>();
  private int playerEntityId = -1;
  private long lastDisposedTick = -1;
  private long lastReceivedTick = -1;
  private long lastTick;
  private long createTick = -1;
  private boolean startedup = false;

  public Engine(Config config, Consumer<FutureOutputMessage> messageOutput,
                int tickRemainder) throws IOException {
    logger.debug("Created engine");
    this.config = config;
    this.messageOutput = messageOutput;
    terrainManager = new TerrainManager();
    terrainImageManager = new TerrainImageManager(terrainManager::getChunk);
    terrainManager.setChunkUpdateConsumer(terrainImageManager::update);
    spriteManager = new SpriteManager();
    spriteManager.loadImages();
    long tick = tickRemainder;
    lastTick = tick;
    createTick = tick;
  }

  private static float getPlayerMovementAngle(TickInput tickInput) {
    float angle = (float) (Math.PI / 4);
    if (tickInput.isMoveRight()) {
      if (tickInput.isMoveUp()) {
        return angle * 1;
      }
      if (tickInput.isMoveDown()) {
        return angle * 7;
      }
      return angle * 0;
    }
    if (tickInput.isMoveLeft()) {
      if (tickInput.isMoveUp()) {
        return angle * 3;
      }
      if (tickInput.isMoveDown()) {
        return angle * 5;
      }
      return angle * 4;
    }
    {
      if (tickInput.isMoveUp()) {
        return angle * 2;
      }
      if (tickInput.isMoveDown()) {
        return angle * 6;
      }
      return Float.NaN;
    }
  }

  public long getTick() {
    return lastTick;
  }

  public void processNextTick(List<Or<Pair.DoubleDouble, Pair.IntBoolean>> eventList) {
    tickInputManager.addInput(lastTick, eventList);
    // redo logic until current tick
    processBufferedMessages();

    logger.trace("Creating tick after {}", lastTick);
    // create next tick
    updateTickAfter(lastTick);
  }

  public void endStartup(long deltaTick) {
    tickInputManager.addInput(lastTick, Collections.emptyList());
    processBufferedMessages();
    long newTick = createTick + deltaTick;
    logger.trace("Creating ticks until {} and finishing engine startup", newTick);
    while (lastTick < newTick) {
      updateTickAfter(lastTick);
    }
    // lastTick = newTick now
    startedup = true;
  }

  public boolean render(Drawer drawer, float alpha, Point.Double mousePosition) {
    if (!startedup) {
      return false;
    }
    TerrainPoint playerPoint;
    if (entityManager.getMapPoint(lastTick - 1, playerEntityId) == null) {
      playerPoint = entityManager.getMapPoint(lastTick, playerEntityId).toTerrainPoint();
    } else {
      playerPoint = entityManager.getMapPoint(lastTick - 1, playerEntityId).toTerrainPoint();
    }
    int xOffset = (int) (MOUSE_SCREEN_MOVING * (mousePosition.x - Ui.getWidth() / 2.0));
    int yOffset = (int) (MOUSE_SCREEN_MOVING * (mousePosition.y - Ui.getHeight() / 2.0));
    TerrainPoint drawCenterPoint = playerPoint.getTranslated(xOffset, yOffset);
    int minChunkX = drawCenterPoint.getChunkX() - (Ui.getWidth() / 2 + 255 - 1 - drawCenterPoint.getBlockX()) / 256 - 1;
    int maxChunkX = drawCenterPoint.getChunkX() + (Ui.getWidth() / 2 - 1 + drawCenterPoint.getBlockX()) / 256 + 1;
    int minChunkY = drawCenterPoint.getChunkY() - (Ui.getHeight() / 2 + 255 - 1 - drawCenterPoint.getBlockY()) / 256 - 1;
    int maxChunkY = drawCenterPoint.getChunkY() + (Ui.getHeight() / 2 - 1 + drawCenterPoint.getBlockY()) / 256 + 1;
    for (int chunkX = minChunkX; chunkX <= maxChunkX; chunkX++) {
      for (int chunkY = minChunkY; chunkY <= maxChunkY; chunkY++) {
        int xMinPos = (chunkX - drawCenterPoint.getChunkX()) * 256 + Ui.getWidth() / 2 - drawCenterPoint.getBlockX();
        int yMinPos = (chunkY - drawCenterPoint.getChunkY()) * 256 + Ui.getHeight() / 2 - drawCenterPoint.getBlockY();
        drawer.pushTranslate(xMinPos, yMinPos);
        terrainImageManager.drawChunk(drawer, new ChunkPoint(chunkX, chunkY));
        drawer.popTranslate();
      }
    }

    final int maxDrawDistance = 50;
    int minXDraw = -maxDrawDistance;
    int maxXDraw = Ui.getWidth() + maxDrawDistance;
    int minYDraw = -maxDrawDistance;
    int maxYDraw = Ui.getHeight() + maxDrawDistance;

    int xOffsetTotal = drawCenterPoint.getX() - Ui.getWidth() / 2;
    int yOffsetTotal = drawCenterPoint.getY() - Ui.getHeight() / 2;

    IntConsumer draw =
            entityId -> {
              TerrainPoint terrainPoint = entityManager.getMapPoint(lastTick, entityId).toTerrainPoint().getTranslated(-xOffsetTotal, -yOffsetTotal);
              if (terrainPoint.getX() < minXDraw || terrainPoint.getX() > maxXDraw) {
                return;
              }
              if (terrainPoint.getY() < minYDraw || terrainPoint.getY() > maxYDraw) {
                return;
              }
              MapPoint oldPoint = entityManager.getMapPoint(lastTick - 1, entityId);
              if (oldPoint != null) {
                TerrainPoint oldterrainPoint = oldPoint.toTerrainPoint().getTranslated(-xOffsetTotal, -yOffsetTotal);
                int deltaX = terrainPoint.getX() - oldterrainPoint.getX();
                int deltaY = terrainPoint.getY() - oldterrainPoint.getY();
                if (deltaX >= -MAX_INTERPOLATION_DISTANCE && deltaX <= MAX_INTERPOLATION_DISTANCE && deltaY >= -MAX_INTERPOLATION_DISTANCE
                        && deltaY <= MAX_INTERPOLATION_DISTANCE) {
                  // old and new are close enough, do lerp
                  terrainPoint = oldterrainPoint.getTranslated((int) alpha * deltaX, (int) alpha * deltaY);
                }
              }
              SpriteType spriteType = entityManager.getSpriteType(lastTick, entityId);
              long spriteTime = entityManager.getSpriteTime(lastTick, entityId);
              spriteManager.drawSpriteType(drawer, spriteType, spriteTime, terrainPoint.getX(), terrainPoint.getY());
            };

    entityManager.getEntitiesStream(lastTick).filter(id -> entityManager.getEntityType(lastTick, id) == EntityType.PLAYER).forEach(draw);
    entityManager.getEntitiesStream(lastTick).filter(id -> entityManager.getEntityType(lastTick, id) == EntityType.BALL).forEach(draw);

    float chatLineHeight = drawer.getLineHeight(Font.COMIC, 12);
    float yChatPosition = chatLineHeight;
    for (String line : chatManager.getMessages()) {
      drawer.drawText(10, yChatPosition, line, Font.COMIC, 12, false, false);
      yChatPosition -= chatLineHeight;
    }

    return true;
  }

  public void freeTime() {
    // we've got some time to dipose the old ticks
    long disposeTick = lastReceivedTick - config.getMaxRollbackTicks();
    terrainManager.disposeUntilTick(disposeTick);
    entityManager.disposeUntilTick(disposeTick);
    playerManager.disposeUntilTick(disposeTick);
    tickInputManager.disposeUntilTick(disposeTick);
    lastDisposedTick = disposeTick;
  }

  private void updateTickAfter(long oldTick) {
    if (oldTick <= lastDisposedTick) {
      throw new IllegalStateException("Tick " + oldTick + " has already been disposed");
    }
    if (oldTick == lastTick) {
      lastTick++;
    }
    long tick = oldTick + 1;

    TickInput tickInput = tickInputManager.getInput(oldTick);

    // update cooldowns and lifespans
    entityManager.getEntitiesStream(tick).forEach(entityId -> {
      if (entityManager.getEntityType(tick, entityId) == EntityType.PLAYER) {
        int cooldown = (int) entityManager.getObjectAttribute(tick, entityId, ObjectAttribute.COOLDOWN_ONE);
        if (cooldown > 0) {
          entityManager.setObjectAttribute(tick, entityId, ObjectAttribute.COOLDOWN_ONE, cooldown - 1);
        }
      }
      if (entityManager.getEntityType(tick, entityId) == EntityType.BALL) {
        int ticksLeft = (int) entityManager.getObjectAttribute(tick, entityId, ObjectAttribute.TICKS_LEFT);
        if (ticksLeft > 0) {
          entityManager.setObjectAttribute(tick, entityId, ObjectAttribute.TICKS_LEFT, ticksLeft - 1);
        } else {
          entityManager.destroyEntity(tick, entityId);
        }
      }
    });

    if (startedup) {
      // set player speed
      float playerMovementDirection = getPlayerMovementAngle(tickInput);
      if (Float.isNaN(playerMovementDirection)) {
        entityManager.setSpeedNorm(tick, playerEntityId, 0f);
      } else {
        entityManager.setSpeedNorm(tick, playerEntityId, config.getDefaultPlayerSpeed());
        entityManager.setSpeedAngle(tick, playerEntityId, playerMovementDirection);
      }
    }

    // integrate positions
    entityManager.getEntitiesStream(tick).forEach(entityId -> {
      float speedNorm = entityManager.getSpeedNorm(tick, entityId);
      if (speedNorm == 0f) {
        return;
      }
      float speedAngle = entityManager.getSpeedAngle(tick, entityId);
      float deltaX = speedNorm * MikoMath.cos(speedAngle);
      float deltaY = speedNorm * MikoMath.sin(speedAngle);
      MapPoint mapPoint = entityManager.getMapPoint(tick, entityId);
      MapPoint newMapPoint = mapPoint.getTranslated(deltaX, deltaY);
      // unused for now
      EntityType entityType = entityManager.getEntityType(tick, entityId);
      // check collisions
      // check terrain collisions
      FloatFloat delta = newMapPoint.subtract(mapPoint);
      int width = (int) delta.getFirst();
      int height = (int) delta.getSecond();
      Hitbox hitbox = entityManager.getSpriteType(tick, entityId).getHitbox();
      for (FloatFloat offset : hitbox.getKeyPoints(delta.getFirst(), delta.getSecond())) {
        TerrainPoint start = mapPoint.getTranslated(offset.getFirst(), offset.getSecond()).toTerrainPoint();
        // brensenham algorithm
        // taken from http://tech-algorithm.com/articles/drawing-line-using-bresenham-algorithm/
        int x = start.getX();
        int y = start.getY();
        int dx1 = 0, dy1 = 0, dx2 = 0, dy2 = 0;
        if (width < 0) {
          dx1 = -1;
        } else if (width > 0) {
          dx1 = 1;
        }
        if (height < 0) {
          dy1 = -1;
        } else if (height > 0) {
          dy1 = 1;
        }
        if (width < 0) {
          dx2 = -1;
        } else if (width > 0) {
          dx2 = 1;
        }
        int longest = Math.abs(width);
        int shortest = Math.abs(height);
        if (longest <= shortest) {
          int temp = longest;
          longest = shortest;
          shortest = temp;
          if (height < 0) {
            dy2 = -1;
          } else if (height > 0) {
            dy2 = 1;
          }
          dx2 = 0;
        }
        int numerator = longest >> 1;
        for (int i = 0; i <= longest; i++) {
          // collision checking with x and y
          TerrainPoint point = new TerrainPoint(x, y);
          TerrainType terrainType = terrainManager.getChunk(tick, point.getChunkPoint()).getBlock(point.getBlockPoint());
          if (terrainType == TerrainType.BLACK_WALL) {
            if (entityType == EntityType.BALL) {
              entityManager.destroyEntity(tick, entityId);
            }
            return;
          }
          numerator += shortest;
          if (!(numerator < longest)) {
            numerator -= longest;
            x += dx1;
            y += dy1;
          } else {
            x += dx2;
            y += dy2;
          }
        }
      }
      // check entity collisions
      // only check balls collisions against other entities
      // only check moving entities against other entities
      if (entityType == EntityType.BALL) {
        for (int otherId : entityManager.getEntities(tick)) {
          int senderId = (int) entityManager.getObjectAttribute(tick, entityId, ObjectAttribute.SENDER);
          if (senderId == otherId) {
            continue;
          }
          if (entityManager.getEntityType(tick, otherId) != EntityType.PLAYER) {
            continue;
          }
          // check hitbox after id checks because it may be more expensive
          Hitbox otherHitbox = entityManager.getSpriteType(tick, otherId).getHitbox();
          MapPoint otherPosition = entityManager.getMapPoint(tick, otherId);
          if (!Hitbox.collide(hitbox, newMapPoint, otherHitbox, otherPosition)) {
            continue;
          }
          // player-ball collision
          entityManager.destroyEntity(tick, entityId);
          int playerHp = (int) entityManager.getObjectAttribute(tick, otherId, ObjectAttribute.HEALTH);
          entityManager.setObjectAttribute(tick, otherId, ObjectAttribute.HEALTH, playerHp - 1);
          return;
        }
      }
      entityManager.setMapPoint(tick, entityId, newMapPoint);
    });

    // process actions
    float ballSendAngle = tickInput.getBallSendRequest();
    if (tick == lastTick && !Float.isNaN(ballSendAngle)) {
      int cooldown = (int) entityManager.getObjectAttribute(oldTick, playerEntityId, ObjectAttribute.COOLDOWN_ONE);
      if (cooldown <= 0) {
        entityManager.setObjectAttribute(tick, playerEntityId, ObjectAttribute.COOLDOWN_ONE, config.getPlayerBallCooldown());
        // create a ball
        int ballId = entityManager.getAndUseTemporaryId();
        MapPoint position = entityManager.getMapPoint(tick, playerEntityId);
        EntityDataUpdate ball =
                new Builder(ballId).entityType(EntityType.BALL).position(position).speedAngle(ballSendAngle)
                        .speedNorm(config.getDefaultBallSpeed()).spriteType(SpriteType.BALL).objectAttribute(ObjectAttribute.SENDER, playerEntityId)
                        .objectAttribute(ObjectAttribute.TICKS_LEFT, config.getDefaultBallLifespan()).build();
        entityManager.createEntity(tick, ball);
        System.out.println("createmoi " + tick + " " + position.getX() + " " + position.getY() + " " + ballSendAngle + " " + config.getDefaultBallSpeed());
        // notify server of ball send
        messageOutput.accept(OutputMessageFactory.action(tick, new Action(ActionType.SEND_BALL, new Pair<>(ballSendAngle, ballId))));
      }
    }

    if (startedup && tick == lastTick) {
      // send client position
      Set<EntityUpdateType> updateTypes = EnumSet.of(EntityUpdateType.POSITION, EntityUpdateType.SPEED_ANGLE, EntityUpdateType.SPEED_NORM);
      EntityDataUpdate playerUpdate = entityManager.generateDataUpdate(tick, playerEntityId, updateTypes, null);
      messageOutput.accept(OutputMessageFactory.entityUpdate(tick, playerUpdate));
    }

    if (tick == lastTick) {
      terrainImageManager.updateTick(lastTick);
    }

    entityManager.getEntitiesStream(tick).forEach(entityId -> {
      if (entityManager.getEntityType(tick, entityId) != EntityType.BALL) {
        return;
      }
      System.out.println("etatmoi " + tick + " " + entityId + " " + entityManager.getMapPoint(tick, entityId).getX() + " " + entityManager.getMapPoint(tick, entityId).getY() + " " + entityManager.getSpeedAngle(tick, entityId) + " " + entityManager.getSpeedNorm(tick, entityId));
    });
  }

  private void processAction(long tick, int entityId, Action action) {
    if (tick <= lastDisposedTick) {
      throw new IllegalStateException("Tick " + tick + " has already been disposed");
    }
    switch (action.getType()) {
      case SEND_BALL:
        // ignorer, la balle sera ajoutée avec un entitycreate/entitiesupdate
        break;
      default:
        throw new IllegalArgumentException("Unknown action type: " + action.getType());
    }
  }

  private String formatMessage(long tick, int entityId, String chatMessage) {
    if (tick <= lastDisposedTick) {
      throw new IllegalStateException("Tick " + tick + " has already been disposed");
    }
    if (playerManager.isPlayer(tick, entityId)) {
      return playerManager.getPlayerName(tick, entityId) + ": " + chatMessage;
    } else {
      return "Entité inconnue: " + chatMessage;
    }
  }

  private void processBufferedMessages() {
    messagesBuffer.sort(Comparator.naturalOrder());
    long previousUpdatedTick = -1;
    Iterator<EngineMessage> it = messagesBuffer.iterator();
    while (it.hasNext()) {
      EngineMessage engineMessage = it.next();
      long tick = engineMessage.getTick();
      if (tick <= lastDisposedTick) {
        throw new IllegalStateException("Tick " + tick + " has already been disposed");
      }
      if (tick > lastTick) {
        if (startedup) {
          logger.warn("Tick " + tick + " has not been created yet");
          //throw new IllegalStateException("Tick " + tick + " has not been created yet");
          break;
        }
        previousUpdatedTick = lastTick;
      }
      if (previousUpdatedTick != -1) {
        // redo logic until this tick
        for (long updateTick = previousUpdatedTick; updateTick < tick; updateTick++) {
          updateTickAfter(updateTick);
        }
      }
      previousUpdatedTick = tick;
      processMessage(engineMessage);
      it.remove();
    }
    if (previousUpdatedTick != -1) {
      lastReceivedTick = previousUpdatedTick > lastReceivedTick ? previousUpdatedTick : lastReceivedTick;
      // update to current tick
      for (long updateTick = previousUpdatedTick; updateTick < lastTick; updateTick++) {
        updateTickAfter(updateTick);
      }
    }
  }

  private void processMessage(EngineMessage engineMessage) {
    long tick = engineMessage.getTick();
    Object[] data = engineMessage.getData();
    switch (engineMessage.getType()) {
      case PLAYER_JOINED:
        int entityIdJoined = (int) data[0];
        String username = (String) data[1];
        playerManager.addPlayer(tick, entityIdJoined, username);
        logger.debug("Player joined with entityId {} as {} on tick {}", entityIdJoined, username, tick);
        break;
      case PLAYER_LEFT:
        int entityIdLeft = (int) data[0];
        playerManager.removePlayer(tick, (int) data[0]);
        logger.debug("Player left with entityId {} on tick {}", entityIdLeft, tick);
        break;
      case CHUNKS_UPDATE:
        @SuppressWarnings("unchecked")
        List<Pair<ChunkPoint, Chunk>> chunks = (List<Pair<ChunkPoint, Chunk>>) data[0];
        for (Pair<ChunkPoint, Chunk> chunk : chunks) {
          terrainManager.setChunk(tick, chunk.getFirst(), chunk.getSecond());
        }
        logger.debug("Chunks update on tick {}, size: {}", tick, chunks.size());
        break;
      case ACTIONS_DONE:
        @SuppressWarnings("unchecked")
        List<Int<Action>> actionList = (List<Int<Action>>) data[0];
        for (Int<Action> actionPair : actionList) {
          processAction(tick, actionPair.getFirst(), actionPair.getSecond());
        }
        logger.debug("Actions received on tick {}", tick);
        break;
      case ENTITY_DESTROY:
        int entityId = (int) data[0];
        entityManager.destroyEntity(tick, entityId);
        logger.debug("Entity destroyed with entityId {} on tick {}", entityId, tick);
        break;
      case ENTITY_CREATE:
        EntityDataUpdate entityDataCreate = (EntityDataUpdate) data[0];
        entityManager.createEntity(tick, entityDataCreate);
        System.out.println("create " + tick + " " + entityDataCreate.getEntityId() + " " + entityDataCreate.getPosition().getX() + " " + entityDataCreate.getPosition().getY() + " " + entityDataCreate.getSpeedAngle() + " " + entityDataCreate.getSpeedNorm() + " " + entityDataCreate.getEntityType());
        logger.debug("Entity created with entityId {} on tick {}", entityDataCreate.getEntityId(), tick);
        break;
      case ENTITIES_UPDATE:
        @SuppressWarnings("unchecked")
        List<EntityDataUpdate> entityDataUpdateList = (List<EntityDataUpdate>) data[0];
        for (EntityDataUpdate entityDataUpdate : entityDataUpdateList) {
          entityManager.applyUpdate(tick, entityDataUpdate);
          System.out.println("update " + tick + " " + entityDataUpdate.getEntityId() + " " + (entityDataUpdate.hasPosition() ? entityDataUpdate.getPosition().getX() + " " + entityDataUpdate.getPosition().getY() + " " : "") + (entityDataUpdate.hasSpeedAngle() ? entityDataUpdate.getSpeedAngle() : "") + " " + (entityDataUpdate.hasSpeedNorm() ? entityDataUpdate.getSpeedNorm() : ""));
        }
        logger.debug("Entities updated on tick {}", tick);
        break;
      default:
        throw new IllegalArgumentException("Unsupported engine message type : " + engineMessage.getType());
    }
  }

  public void actions(int tickRemainder, List<Int<Action>> actions) {
    messagesBuffer.add(EngineMessage.newActionsMessage(getTick(tickRemainder), actions));
  }

  public void chatReceived(int tickRemainder, int entityIdChat, String chatMessage) {
    String formattedMessage = formatMessage(getTick(tickRemainder), entityIdChat, chatMessage);
    chatManager.addMessage(formattedMessage);
    logger.info("Received chat message from entityId {}", entityIdChat);
  }

  public void chunksUpdate(int tickRemainder, List<Pair<ChunkPoint, Chunk>> chunks) {
    messagesBuffer.add(EngineMessage.newChunksUpdateMessage(getTick(tickRemainder), chunks));
  }

  public void entityIdChange(int oldEntityId, int newEntityId) {
    entityManager.freeAndUpdateTemporaryId(oldEntityId, newEntityId);
    logger.debug("Received entity id change from {} to {}", oldEntityId, newEntityId);
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
    logger.debug("Set player entity id to {}", entityId);
  }

  public void wroteMessage(String message) {
    // TODO call this from miko
    String formattedMessage = formatMessage(lastTick, playerEntityId, message);
    chatManager.addMessage(formattedMessage);
    messageOutput.accept(OutputMessageFactory.chatSend(message));
    logger.info("Sent chat message");
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
    return newTick;
  }
}
