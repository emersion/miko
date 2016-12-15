package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.*;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Map.Entry;
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
  private static final int TEMPORARY_ID_START = 64536;
  private static Logger logger = LogManager.getLogger("miko.engine");
  private boolean[] usedTemporaryIds = new boolean[(1 << 16) - TEMPORARY_ID_START];
  // improves performance on temporaryid generating
  private int circularTemporaryId = usedTemporaryIds.length - 1;
  private Map<Integer, Entity> map = new HashMap<>();
  // cache last entity to improve performance on repeated get and set calls
  private int lastEntityId = -1;
  private Entity lastEntity;

  public void applyUpdate(long tick, EntityDataUpdate entityDataUpdate) {
    getEntity(entityDataUpdate.getEntityId()).applyUpdate(tick, entityDataUpdate);
  }

  public void createEntity(long tick, EntityDataUpdate entityDataUpdate) {
    Entity entity = getEntity(entityDataUpdate.getEntityId());
    if (entity == null) {
      entity = new Entity();
      map.put(entityDataUpdate.getEntityId(), entity);
      lastEntityId = -1;
    }
    entity.applyFullUpdate(tick, entityDataUpdate);
  }

  public void destroyEntity(long tick, int entityId) {
    Entity entity = getEntity(entityId);
    if(entity == null) {
      logger.warn("Tried destroying entity {} at tick {} but it doesn't exist", entityId, tick);
      return;
    }
    entity.disable(tick);
  }

  public EntityDataUpdate generateDataUpdate(long tick, int entityId, Set<EntityUpdateType> types, Set<ObjectAttribute> attributes) {
    return getEntity(entityId).generateUpdate(entityId, tick, types, attributes);
  }

  public EntityType getEntityType(long tick, int entityId) {
    return getEntity(entityId).getEntityType(tick);
  }

  public MapPoint getMapPoint(long tick, int entityId) {
    return getEntity(entityId).getMapPoint(tick);
  }

  public float getSpeedAngle(long tick, int entityId) {
    return getEntity(entityId).getSpeedAngle(tick);
  }

  public float getSpeedNorm(long tick, int entityId) {
    return getEntity(entityId).getSpeedNorm(tick);
  }

  public SpriteType getSpriteType(long tick, int entityId) {
    return getEntity(entityId).getSpriteType(tick);
  }

  public long getSpriteTime(long tick, int entityId) {
    return getEntity(entityId).getSpriteTime(tick);
  }

  public Object getObjectAttribute(long tick, int entityId, ObjectAttribute type) {
    return getEntity(entityId).getObjectAttribute(tick, type);
  }

  public boolean isEnabled(long tick, int entityId) {
    return getEntity(entityId).isEnabled(tick);
  }

  public void setMapPoint(long tick, int entityId, MapPoint mapPoint) {
    getEntity(entityId).setMapPoint(tick, mapPoint);
  }

  public void setTerrainPoint(long tick, int entityId, TerrainPoint terrainPoint) {
    getEntity(entityId).setTerrainPoint(tick, terrainPoint);
  }

  public void setSpeedAngle(long tick, int entityId, float speedAngle) {
    getEntity(entityId).setSpeedAngle(tick, speedAngle);
  }

  public void setSpeedNorm(long tick, int entityId, float speedNorm) {
    getEntity(entityId).setSpeedNorm(tick, speedNorm);
  }

  public void setSpriteType(long tick, int entityId, SpriteType spriteType) {
    getEntity(entityId).setSpriteType(tick, spriteType);
  }

  public void setObjectAttribute(long tick, int entityId, ObjectAttribute attribute, Object value) {
    getEntity(entityId).setObjectAttribute(tick, attribute, value);
  }

  /**
   * Retourne et alloue un entityId temporaire à l'engine. Utilisé pour ajouter des objets sans
   * connaître leur id.
   *
   * @return Un entityId temporaire.
   * @see #freeAndUpdateTemporaryId(int, int)
   */
  public int getAndUseTemporaryId() {
    int start = circularTemporaryId + 1;
    int end = circularTemporaryId + usedTemporaryIds.length;
    for (int i = start; i <= end; i++) {
      int id = i % usedTemporaryIds.length;
      if (!usedTemporaryIds[id]) {
        usedTemporaryIds[id] = true;
        circularTemporaryId = id;
        return id + TEMPORARY_ID_START;
      }
    }
    throw new RuntimeException("Ran out of temporary entities id");
  }

  /**
   * Désalloue un entityId temporaire et change l'entityId de l'entité éventuelle occupant cet id
   * par un nouvel id (permanent).
   * <p>
   * S'il n'y avait aucune entité à l'entityId temporaire spécifié, le nouvel id sera ignoré.
   *
   * @param oldEntityId L'entityId de l'entité temporaire alloué précédemment.
   * @param newEntityId Le nouvel entityId, permanent, à associer à l'entité.
   */
  public void freeAndUpdateTemporaryId(int oldEntityId, int newEntityId) {
    int offsetId = oldEntityId - TEMPORARY_ID_START;
    if (offsetId < 0) {
      throw new IllegalArgumentException("oldEntityId must be a temporary id");
    }
    if (!usedTemporaryIds[offsetId]) {
      return;
    }
    usedTemporaryIds[offsetId] = false;
    Entity entity = map.remove(oldEntityId);
    if (entity == null) {
      return;
    }
    if (newEntityId >= TEMPORARY_ID_START) {
      throw new IllegalArgumentException("newEntityId must be a permanent (normal) id");
    }
    map.put(newEntityId, entity);
    lastEntityId = -1;
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
    Set<Entry<Integer, Entity>> entrySet = map.entrySet();
    Iterator<Entry<Integer, Entity>> entryIterator = entrySet.iterator();
    while (entryIterator.hasNext()) {
      Entry<Integer, Entity> entry = entryIterator.next();
      Entity entity = entry.getValue();
      entity.disposeUntilTick(tick);
      if (entity.isEmpty()) {
        entryIterator.remove();
      }
    }
  }

  public IntStream getEntitiesStream(long tick) {
    return map.entrySet().stream().filter(e -> e.getValue().isEnabled(tick)).mapToInt(e -> e.getKey());
  }

  public Iterable<Integer> getEntities(long tick) {
    return getEntitiesStream(tick).boxed()::iterator;
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
