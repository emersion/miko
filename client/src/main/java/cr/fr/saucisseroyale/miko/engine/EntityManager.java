package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.ObjectAttribute;
import cr.fr.saucisseroyale.miko.protocol.Sprite;

import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.stream.IntStream;

/**
 * Un gestionnaire des entités du jeu, stockant toutes les entités à tous les ticks, avec le
 * principe de {@link Snapshots}.
 *
 * @see Entity
 * @see Snapshots
 */
class EntityManager {

  private Map<Integer, Entity> map;

  // cache last entity to improve performance on repeated get and set calls
  private int lastEntityId = -1;
  private Entity lastEntity;

  public void applyUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    getEntity(entityDataUpdate.getEntityId()).applyUpdate(tick, entityDataUpdate);
  }

  public void createEntity(long tick, EntityDataUpdate entityDataUpdate) {
    getEntity(entityDataUpdate.getEntityId()).applyFullUpdate(tick, entityDataUpdate);
  }

  public void destroyEntity(long tick, int entityId) {
    getEntity(entityId).disable(tick);
  }

  public BlockPoint getBlockPoint(long tick, int entityId) {
    return getEntity(entityId).getBlockPoint(tick);
  }

  public ChunkPoint getChunkPoint(long tick, int entityId) {
    return getEntity(entityId).getChunkPoint(tick);
  }

  public float getSpeedAngle(long tick, int entityId) {
    return getEntity(entityId).getSpeedAngle(tick);
  }

  public float getSpeedNorm(long tick, int entityId) {
    return getEntity(entityId).getSpeedNorm(tick);
  }

  public Sprite getSprite(int entityId) {
    return getEntity(entityId).getSprite();
  }

  public long getSpriteTime(long tick, int entityId) {
    return getEntity(entityId).getSpriteTime(tick);
  }

  public Set<ObjectAttribute> getObjectAttributes(long tick, int entityId) {
    return getEntity(entityId).getObjectAttributes(tick);
  }

  public boolean isEnabled(long tick, int entityId) {
    return getEntity(entityId).isEnabled(tick);
  }

  public void setBlockPoint(long tick, int entityId, BlockPoint blockPoint) {
    getEntity(entityId).setBlockPoint(tick, blockPoint);
  }

  public void setChunkPoint(long tick, int entityId, ChunkPoint chunkPoint) {
    getEntity(entityId).setChunkPoint(tick, chunkPoint);
  }

  public void setSpeedAngle(long tick, int entityId, float speedAngle) {
    getEntity(entityId).setSpeedAngle(tick, speedAngle);
  }

  public void setSpeedNorm(long tick, int entityId, float speedNorm) {
    getEntity(entityId).setSpeedNorm(tick, speedNorm);
  }

  public void setSprite(long tick, int entityId, Sprite sprite) {
    getEntity(entityId).setSprite(tick, sprite);
  }

  public void setObjectAttributes(long tick, int entityId, Set<ObjectAttribute> objectAttributes) {
    getEntity(entityId).setObjectAttributes(tick, objectAttributes);
  }

  /**
   * Indique au gestionnaire d'entités que les entités appartenant à des ticks avant ou égaux au
   * tick spécifié ne seront plus jamais demandées et peuvent être supprimées.
   *
   * @param tick Le tick (inclus) jusqu'auquel les entités ne seront plus demandées.
   */
  public void disposeUntilTick(long tick) {
    // NOTE on pourrait ignorer ce call la plupart du temps pour éviter de parcourir toute la map à
    // chaque fois ; il faudrait sans doute modifier la capacité standard
    Set<Map.Entry<Integer, Entity>> entrySet = map.entrySet();
    Iterator<Map.Entry<Integer, Entity>> entryIterator = entrySet.iterator();
    while (entryIterator.hasNext()) {
      Map.Entry<Integer, Entity> entry = entryIterator.next();
      Entity entity = entry.getValue();
      entity.disposeUntilTick(tick);
      if (entity.isEmpty()) {
        entryIterator.remove();
      }
    }
  }

  public IntStream getEntities(long tick) {
    return map.entrySet().stream().filter((e) -> e.getValue().isEnabled(tick)).mapToInt((e) -> e.getKey());
  }

  private Entity getEntity(int entityId) {
    if (entityId == lastEntityId) {
      return lastEntity;
    }
    lastEntityId = entityId;
    lastEntity = map.get(entityId);
    return lastEntity;
  }

}
