package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.network.FutureOutputMessage;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.Config;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.EntityType;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.protocol.TerrainPoint;
import cr.fr.saucisseroyale.miko.protocol.TerrainType;
import cr.fr.saucisseroyale.miko.util.MikoMath;
import cr.fr.saucisseroyale.miko.util.Pair;
import cr.fr.saucisseroyale.miko.util.Triplet;

import java.awt.Graphics2D;
import java.awt.GraphicsConfiguration;
import java.awt.Point;
import java.awt.image.BufferedImage;
import java.awt.image.DataBuffer;
import java.awt.image.DataBufferInt;
import java.awt.image.DirectColorModel;
import java.awt.image.Raster;
import java.awt.image.WritableRaster;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.function.Consumer;
import java.util.function.IntConsumer;

/**
 * L'engine du jeu Miko, écoutant les inputs du serveur et du client, et pouvant être affiché.
 *
 */
public class Engine {

  private static final int TICK_DIVIDER = 1 << 16;
  // sets how much the screen moves in response to mouse moves
  // MUST be positive
  private static final float MOUSE_SCREEN_MOVING = 0.2f;

  private Config config;
  private int screenWidth;
  private int screenHeight;

  private BufferedImage terrainImage;
  private int[] terrainImageData;

  private TerrainManager terrainManager = new TerrainManager();
  private EntityManager entityManager = new EntityManager();
  private PlayerManager playerManager = new PlayerManager();
  private ChatManager chatManager = new ChatManager();
  private TickInputManager tickInputManager = new TickInputManager();
  private SpriteManager spriteManager;

  private List<EngineMessage> messagesBuffer;
  private final Consumer<FutureOutputMessage> messageOutput;

  private int playerEntityId = -1;

  private long lastDisposedTick = Long.MIN_VALUE;
  private long lastTick;

  public Engine(Config config, GraphicsConfiguration graphicsConfiguration, Consumer<FutureOutputMessage> messageOutput, int width, int height,
      int tickRemainder) throws IOException {
    this.config = config;
    this.messageOutput = messageOutput;
    screenWidth = width;
    screenHeight = height;
    terrainImageData = new int[screenWidth * screenHeight];
    DataBuffer buffer = new DataBufferInt(terrainImageData, terrainImageData.length);
    WritableRaster raster = Raster.createPackedRaster(buffer, screenWidth, screenHeight, screenWidth, new int[] {0xFF0000, 0xFF00, 0xFF}, null);
    terrainImage = new BufferedImage(new DirectColorModel(24, 0xFF0000, 0xFF00, 0xFF), raster, false, null);
    spriteManager = new SpriteManager(graphicsConfiguration);
    spriteManager.loadImages();
    long tick = tickRemainder;
    lastTick = tick;
    messagesBuffer = new ArrayList<>();
  }

  public void processNextTick(List<Triplet<Boolean, Integer, Point>> eventList) {
    tickInputManager.addInput(lastTick, eventList);
    // redo logic until current tick
    processBufferedMessages();
    // create next tick
    updateTickAfter(lastTick);
  }

  public void render(Graphics2D graphics, float alpha, Point mousePosition) {
    TerrainPoint playerPoint = entityManager.getMapPoint(lastTick, playerEntityId).toTerrainPoint();
    int xOffset = (int) (MOUSE_SCREEN_MOVING * mousePosition.x);
    int yOffset = (int) (MOUSE_SCREEN_MOVING * mousePosition.y);
    TerrainPoint drawCenterPoint = playerPoint.getTranslated(xOffset, yOffset);

    int minChunkX = drawCenterPoint.getChunkX() - (screenWidth / 2 - drawCenterPoint.getBlockX()) / 256;
    int maxChunkX = drawCenterPoint.getChunkX() + (screenWidth / 2 + drawCenterPoint.getBlockX()) / 256;
    int minChunkY = drawCenterPoint.getChunkY() - (screenHeight / 2 - drawCenterPoint.getBlockY()) / 256;
    int maxChunkY = drawCenterPoint.getChunkY() + (screenHeight / 2 + drawCenterPoint.getBlockY()) / 256;
    for (int chunkX = minChunkX; chunkX <= maxChunkX; chunkX++) {
      for (int chunkY = minChunkY; chunkY <= maxChunkY; chunkY++) {
        Chunk chunk = terrainManager.getChunk(lastTick, new ChunkPoint(chunkX, chunkY));
        int defaultColor = chunk.getDefaultType().getColor();
        int xMinPos = (chunkX - drawCenterPoint.getChunkX()) * 256 + screenWidth / 2 - drawCenterPoint.getBlockX();
        int yMinPos = (chunkY - drawCenterPoint.getChunkY()) * 256 + screenHeight / 2 - drawCenterPoint.getBlockY();
        int xMinOffset = xMinPos < 0 ? -xMinPos : 0;
        int xMaxOffset = xMinPos + 255 >= screenWidth ? screenWidth - xMinPos - 1 : 255;
        int yMinOffset = yMinPos < 0 ? -yMinPos : 0;
        int yMaxOffset = yMinPos + 255 >= screenHeight ? screenHeight - yMinPos - 1 : 255;
        int startPos = xMinPos + (screenHeight - (yMinPos + yMaxOffset)) * screenWidth;
        int endPos = xMinPos + (screenHeight - (yMinPos + yMinOffset)) * screenWidth;
        for (int pos = startPos; pos <= endPos; pos += screenWidth) {
          for (int i = xMinOffset; i <= xMaxOffset; i++) {
            terrainImageData[pos + i] = defaultColor;
          }
        }
        for (Block block : chunk.getDefinedBlocks()) {
          if (block.getX() < xMinOffset || block.getX() > xMaxOffset || block.getY() < yMinOffset || block.getY() > yMaxOffset) {
            continue;
          }
          int pos = xMinPos + block.getX() + (screenHeight - (yMinPos + block.getY())) * screenWidth;
          terrainImageData[pos] = block.getType().getColor();
        }
      }
    }
    graphics.drawImage(terrainImage, 0, screenHeight, null);

    final int maxDrawDistance = 50;
    int minXDraw = drawCenterPoint.getX() - screenWidth / 2 - maxDrawDistance;
    int maxXDraw = drawCenterPoint.getX() + screenWidth / 2 + maxDrawDistance;
    int minYDraw = drawCenterPoint.getY() - screenHeight / 2 - maxDrawDistance;
    int maxYDraw = drawCenterPoint.getY() + screenHeight / 2 + maxDrawDistance;

    int xOffsetTotal = drawCenterPoint.getX();
    int yOffsetTotal = drawCenterPoint.getY();

    // TODO interpolation linéaire
    IntConsumer draw = (entityId) -> {
      TerrainPoint terrainPoint = entityManager.getMapPoint(lastTick, entityId).toTerrainPoint().getTranslated(xOffsetTotal, yOffsetTotal);
      if (terrainPoint.getX() < minXDraw || terrainPoint.getX() > maxXDraw) {
        return;
      }
      if (terrainPoint.getY() < minYDraw || terrainPoint.getY() > maxYDraw) {
        return;
      }
      SpriteType spriteType = entityManager.getSpriteType(lastTick, entityId);
      long spriteTime = entityManager.getSpriteTime(lastTick, entityId);
      spriteManager.drawSpriteType(graphics, spriteType, spriteTime, terrainPoint.getX(), terrainPoint.getY());
    };
    entityManager.getEntitiesStream(lastTick).filter((id) -> entityManager.getEntityType(lastTick, id) == EntityType.PLAYER).forEach(draw);
    entityManager.getEntitiesStream(lastTick).filter((id) -> entityManager.getEntityType(lastTick, id) == EntityType.BALL).forEach(draw);

    // TODO draw chat
  }

  public void freeTime() {
    // we've got some time to dipose the old ticks
    long disposeTick = lastTick - (config.getMaxRollbackTicks() + 1); // +1 to be safe
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
    entityManager.getEntitiesStream(tick).forEach((entityId) -> {
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

    // set player speed
    float playerMovementDirection = getPlayerMovementAngle(tickInput);
    if (playerMovementDirection == Float.NaN) {
      entityManager.setSpeedNorm(tick, playerEntityId, 0f);
    } else {
      entityManager.setSpeedNorm(tick, playerEntityId, config.getDefaultPlayerSpeed());
      entityManager.setSpeedAngle(tick, playerEntityId, playerMovementDirection);
    }

    // integrate positions
    entityManager.getEntitiesStream(tick).forEach((entityId) -> {
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
      Pair<Float, Float> delta = newMapPoint.subtract(mapPoint);
      int width = delta.getFirst().intValue();
      int height = delta.getSecond().intValue();
      for (Pair<Float, Float> offset : entityManager.getSpriteType(tick, entityId).getHitbox().getKeyPoints(delta.getFirst(), delta.getSecond())) {
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
              return;
            }
            if (entityType == EntityType.PLAYER) {
              entityManager.setSpeedNorm(tick, entityId, 0f);
            }
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
    for (Float ballSendAngle : tickInput.getBallSendRequests()) {
      int cooldown = (int) entityManager.getObjectAttribute(oldTick, playerEntityId, ObjectAttribute.COOLDOWN_ONE);
      if (cooldown > 0) {
        break;
      }
      entityManager.setObjectAttribute(tick, playerEntityId, ObjectAttribute.COOLDOWN_ONE, config.getPlayerBallCooldown());
      // create a ball
      int ballId = entityManager.getAndUseTemporaryId();
      MapPoint position = entityManager.getMapPoint(tick, playerEntityId);
      EntityDataUpdate ball =
          new EntityDataUpdate.Builder(ballId).entityType(EntityType.BALL).position(position).speedAngle(ballSendAngle)
          .speedNorm(config.getDefaultBallSpeed()).spriteType(SpriteType.BALL).objectAttribute(ObjectAttribute.SENDER, playerEntityId)
          .objectAttribute(ObjectAttribute.TICKS_LEFT, config.getDefaultBallLifespan()).build();
      entityManager.createEntity(tick, ball);
    }

  }

  private void processAction(long tick, int entityId, Action action) {
    if (tick <= lastDisposedTick) {
      throw new IllegalStateException("Tick " + tick + " has already been disposed");
    }
    switch (action.getType()) {
      case SEND_BALL:
        // ignorer, la balle sera ajoutée avec un entitycreate/entitiesupdate
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
          terrainManager.setChunk(tick, (ChunkPoint) data[0], (Chunk) data[1]);
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
}
