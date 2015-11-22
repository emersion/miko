package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.EntityType;
import cr.fr.saucisseroyale.miko.protocol.EntityUpdateType;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.protocol.TerrainPoint;

import java.util.Map;
import java.util.Set;

/**
 * Une entité existant sur plusieurs ticks, supportant des mises à jour partielles.
 *
 */
class Entity {

  // TODO vérifier si c'est plus rapide avec arraysnaphosts qu'avec linkedsnapshots

  private final Snapshots<MapPoint> mapPoints;
  private final Snapshots<Float> speedAngles;
  private final Snapshots<Float> speedNorms;
  private final Snapshots<EntityType> entityTypes;
  private final SnapshotsMap<ObjectAttribute, Object> objectAttributes;
  private final Snapshots<Boolean> enableStatus;
  private final Snapshots<SpriteType> spriteTypes;

  private long tickSwitchedSprite;

  public Entity() {
    mapPoints = new Snapshots<>();
    speedAngles = new Snapshots<>();
    speedNorms = new Snapshots<>();
    entityTypes = new Snapshots<>();
    objectAttributes = new SnapshotsMap<>();
    enableStatus = new Snapshots<>();
    spriteTypes = new Snapshots<>();
  }

  public void applyUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    if (entityDataUpdate.hasPosition()) {
      setMapPoint(tick, entityDataUpdate.getPosition());
    }
    if (entityDataUpdate.hasSpeedAngle()) {
      setSpeedAngle(tick, entityDataUpdate.getSpeedAngle());
    }
    if (entityDataUpdate.hasSpeedNorm()) {
      setSpeedNorm(tick, entityDataUpdate.getSpeedNorm());
    }
    if (entityDataUpdate.hasEntityType()) {
      setEntityType(tick, entityDataUpdate.getEntityType());
    }
    if (entityDataUpdate.hasSprite()) {
      setSpriteType(tick, entityDataUpdate.getSpriteType());
    }
    setObjectAttributes(tick, entityDataUpdate.getObjectAttributes());
  }

  public void applyFullUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    enableStatus.setSnapshot(tick, Boolean.TRUE);
    if (entityDataUpdate.hasPosition()) {
      setMapPoint(tick, entityDataUpdate.getPosition());
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
    if (entityDataUpdate.hasEntityType()) {
      setEntityType(tick, entityDataUpdate.getEntityType());
    } else {
      throw new IllegalArgumentException("Illegal full data update: entityType not set");
    }
    if (entityDataUpdate.hasSprite()) {
      setSpriteType(tick, entityDataUpdate.getSpriteType());
    } else {
      throw new IllegalArgumentException("Illegal full data update: spriteType not set");
    }
    setObjectAttributes(tick, entityDataUpdate.getObjectAttributes());
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

  public EntityDataUpdate generateUpdate(int entityId, long tick, Set<EntityUpdateType> types, Set<ObjectAttribute> attributeTypes) {
    EntityDataUpdate.Builder builder = new EntityDataUpdate.Builder(entityId);
    for (EntityUpdateType type : types) {
      switch (type) {
        case POSITION:
          builder.position(getMapPoint(tick));
          break;
        case SPEED_ANGLE:
          builder.speedAngle(getSpeedAngle(tick));
          break;
        case SPEED_NORM:
          builder.speedNorm(getSpeedNorm(tick));
          break;
        case ENTITY_TYPE:
          builder.entityType(getEntityType(tick));
          break;
        case SPRITE_TYPE:
          builder.spriteType(getSpriteType(tick));
          break;
        case OBJECT_DATA:
          for (ObjectAttribute attributeType : attributeTypes) {
            builder.objectAttribute(attributeType, getObjectAttribute(tick, attributeType));
          }
          break;
        default:
          throw new IllegalArgumentException("attribute " + type + " is not supported");
      }
    }
    return builder.build();
  }

  // no enable check before getting to maximize performance
  // we trust users not to get on disabled entities
  // if this is a problem just add a check in the future

  public EntityType getEntityType(long tick) {
    return entityTypes.getSnapshot(tick);
  }

  public MapPoint getMapPoint(long tick) {
    return mapPoints.getSnapshot(tick);
  }

  public float getSpeedAngle(long tick) {
    return speedAngles.getSnapshot(tick);
  }

  public float getSpeedNorm(long tick) {
    return speedNorms.getSnapshot(tick);
  }

  public Object getObjectAttribute(long tick, ObjectAttribute type) {
    return objectAttributes.getSnapshot(tick, type);
  }

  public SpriteType getSpriteType(long tick) {
    return spriteTypes.getSnapshot(tick);
  }

  public long getSpriteTime(long tick) {
    return tick - tickSwitchedSprite;
  }

  public void setEntityType(long tick, EntityType entityType) {
    entityTypes.setSnapshot(tick, entityType);
  }

  public void setMapPoint(long tick, MapPoint mapPoint) {
    mapPoints.setSnapshot(tick, mapPoint);
  }

  public void setTerrainPoint(long tick, TerrainPoint terrainPoint) {
    mapPoints.setSnapshot(tick, new MapPoint(terrainPoint));
  }

  public void setSpeedAngle(long tick, float speedAngle) {
    speedAngles.setSnapshot(tick, speedAngle);
  }

  public void setSpeedNorm(long tick, float speedNorm) {
    speedNorms.setSnapshot(tick, speedNorm);
  }

  public void setObjectAttribute(long tick, ObjectAttribute attribute, Object value) {
    objectAttributes.setSnapshot(tick, attribute, value);
  }

  private void setObjectAttributes(long tick, Map<ObjectAttribute, Object> newObjectAttributes) {
    for (Map.Entry<ObjectAttribute, Object> newObjetAttribute : newObjectAttributes.entrySet()) {
      setObjectAttribute(tick, newObjetAttribute.getKey(), newObjetAttribute.getValue());
    }
  }

  public void setSpriteType(long tick, SpriteType spriteType) {
    spriteTypes.setSnapshot(tick, spriteType);
    if (tick > tickSwitchedSprite) {
      tickSwitchedSprite = tick;
    }
  }

  public void disposeUntilTick(long tick) {
    mapPoints.disposeUntilTick(tick);
    speedAngles.disposeUntilTick(tick);
    speedNorms.disposeUntilTick(tick);
    enableStatus.disposeUntilTick(tick);
    entityTypes.disposeUntilTick(tick);
    objectAttributes.disposeUntilTick(tick);
    spriteTypes.disposeUntilTick(tick);
  }
}
