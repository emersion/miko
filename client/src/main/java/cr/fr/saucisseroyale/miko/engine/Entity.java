package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.MapPoint;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.Sprite;

import java.util.Set;

/**
 * Une entité existant sur plusieurs ticks, supportant des mises à jour partielles.
 *
 */
class Entity {

  private Snapshots<BlockPoint> blockPoints;
  private Snapshots<ChunkPoint> chunkPoints;
  private Snapshots<Float> speedAngles;
  private Snapshots<Float> speedNorms;
  private Snapshots<Set<ObjectAttribute>> objectAttributeSets;
  private Snapshots<Boolean> enableStatus;

  private long tickSwitchedSprite;
  private Sprite sprite;

  public void applyUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    if (entityDataUpdate.hasPosition()) {
      MapPoint update = entityDataUpdate.getPosition();
      setChunkPoint(tick, new ChunkPoint(update.getChunkX(), update.getChunkY()));
      setBlockPoint(tick, new BlockPoint(update.getBlockX(), update.getBlockY()));
    }
    if (entityDataUpdate.hasSpeedAngle()) {
      setSpeedAngle(tick, entityDataUpdate.getSpeedAngle());
    }
    if (entityDataUpdate.hasSpeedNorm()) {
      setSpeedNorm(tick, entityDataUpdate.getSpeedNorm());
    }
    if (entityDataUpdate.hasSprite()) {
      setSprite(tick, entityDataUpdate.getSprite());
    }
    if (entityDataUpdate.hasObjectAttributes()) {
      setObjectAttributes(tick, entityDataUpdate.getObjectAttributes());
    }
  }

  public void applyFullUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    enableStatus.setSnapshot(tick, Boolean.TRUE);
    if (entityDataUpdate.hasPosition()) {
      MapPoint update = entityDataUpdate.getPosition();
      setChunkPoint(tick, new ChunkPoint(update.getChunkX(), update.getChunkY()));
      setBlockPoint(tick, new BlockPoint(update.getBlockX(), update.getBlockY()));
    } else {
      throw new IllegalArgumentException("Illegal full data update: position not set");
    }
    if (entityDataUpdate.hasSpeedAngle()) {
      setSpeedAngle(tick, entityDataUpdate.getSpeedAngle());
    } else {
      throw new IllegalArgumentException("Illegal full data update: speedangle not set");
    }
    if (entityDataUpdate.hasSpeedNorm()) {
      setSpeedNorm(tick, entityDataUpdate.getSpeedNorm());
    } else {
      throw new IllegalArgumentException("Illegal full data update: speednorm not set");
    }
    if (entityDataUpdate.hasSprite()) {
      setSprite(tick, entityDataUpdate.getSprite());
    } else {
      throw new IllegalArgumentException("Illegal full data update: sprite not set");
    }
    if (entityDataUpdate.hasObjectAttributes()) {
      setObjectAttributes(tick, entityDataUpdate.getObjectAttributes());
    } else {
      throw new IllegalArgumentException("Illegal full data update: objectattributes not set");
    }
  }

  public boolean isEnabled(long tick) {
    Boolean status = enableStatus.getSnapshot(tick);
    if (status == null) {
      return false;
    }
    return status;
  }

  public void disable(long tick) {
    enableStatus.setSnapshot(tick, Boolean.FALSE);
  }

  public boolean isEmpty() {
    if (enableStatus.isEmpty()) {
      return true;
    }
    if (enableStatus.size() == 1 && !isEnabled(Long.MAX_VALUE)) {
      return true;
    }
    return false;
  }

  // no enable check before getting to maximize performance
  // we trust users not to get on disabled entities
  // if this is a problem just add a check in the future

  public BlockPoint getBlockPoint(long tick) {
    return blockPoints.getSnapshot(tick);
  }

  public ChunkPoint getChunkPoint(long tick) {
    return chunkPoints.getSnapshot(tick);
  }

  public float getSpeedAngle(long tick) {
    return speedAngles.getSnapshot(tick);
  }

  public float getSpeedNorm(long tick) {
    return speedNorms.getSnapshot(tick);
  }

  public Set<ObjectAttribute> getObjectAttributes(long tick) {
    return objectAttributeSets.getSnapshot(tick);
  }

  public Sprite getSprite() {
    return sprite;
  }

  public long getSpriteTime(long tick) {
    return tick - tickSwitchedSprite;
  }

  public void setBlockPoint(long tick, BlockPoint blockPoint) {
    blockPoints.setSnapshot(tick, blockPoint);
  }

  public void setChunkPoint(long tick, ChunkPoint chunkPoint) {
    chunkPoints.setSnapshot(tick, chunkPoint);
  }

  public void setSpeedAngle(long tick, float speedAngle) {
    speedAngles.setSnapshot(tick, speedAngle);
  }

  public void setSpeedNorm(long tick, float speedNorm) {
    speedNorms.setSnapshot(tick, speedNorm);
  }

  public void setObjectAttributes(long tick, Set<ObjectAttribute> objectAttributeSet) {
    objectAttributeSets.setSnapshot(tick, objectAttributeSet);
  }

  public void setSprite(long tick, Sprite sprite) {
    this.sprite = sprite;
    tickSwitchedSprite = tick;
  }

  public void disposeUntilTick(long tick) {
    blockPoints.disposeUntilTick(tick);
    chunkPoints.disposeUntilTick(tick);
    speedAngles.disposeUntilTick(tick);
    speedNorms.disposeUntilTick(tick);
    objectAttributeSets.disposeUntilTick(tick);
  }

}
