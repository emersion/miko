package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.engine.MapPoint;

import java.util.Collections;
import java.util.EnumMap;
import java.util.Map;
import java.util.NoSuchElementException;

/**
 * Un objet immutable stockant des valeurs pour mettre à jour une entité. Suit la pattern Builder.
 *
 * @see Builder
 *
 */
public final class EntityDataUpdate {

  /**
   * Un builder pour la classe {@link EntityDataUpdate}. Ne peut être construit qu'une seule fois.
   *
   */
  public static class Builder {
    private final int entityId;
    private final Map<EntityUpdateType, Object> data = new EnumMap<>(EntityUpdateType.class);
    private final Map<ObjectAttribute, Object> objectAttributes = new EnumMap<>(ObjectAttribute.class);
    private boolean built = false;

    /**
     * @param entityId L'id de l'entité à mettre à jour.
     */
    public Builder(int entityId) {
      if (entityId < 0 || entityId >= 1 << 16) {
        throw new IllegalArgumentException("entityId must be between 0 and 65535 inclusive");
      }
      this.entityId = entityId;
    }

    /**
     *
     * @param position La nouvelle position de l'entité.
     * @return self
     */
    public Builder position(TerrainPoint position) {
      return position(new MapPoint(position));
    }

    /**
     *
     * @param position La nouvelle position de l'entité.
     * @return self
     */
    public Builder position(MapPoint position) {
      ensureNotBuilt();
      data.put(EntityUpdateType.POSITION, position);
      return this;
    }

    /**
     *
     * @param speedAngle Le nouvel angle du vecteur vitesse.
     * @return self
     */
    public Builder speedAngle(float speedAngle) {
      ensureNotBuilt();
      data.put(EntityUpdateType.SPEED_ANGLE, speedAngle);
      return this;
    }

    /**
     *
     * @param speedNorm La nouvelle norme du vecteur vitesse.
     * @return self
     */
    public Builder speedNorm(float speedNorm) {
      ensureNotBuilt();
      data.put(EntityUpdateType.SPEED_NORM, speedNorm);
      return this;
    }

    /**
     *
     * @param entityType Le nouveau type d'entité.
     * @return self
     */

    public Builder entityType(EntityType entityType) {
      ensureNotBuilt();
      data.put(EntityUpdateType.ENTITY_TYPE, entityType);
      return this;
    }

    /**
     *
     * @param spriteType Le nouveau type de sprite de l'entité.
     * @return self
     */
    public Builder spriteType(SpriteType spriteType) {
      ensureNotBuilt();
      data.put(EntityUpdateType.SPRITE_TYPE, spriteType);
      return this;
    }

    /**
     *
     * @param type Le type de l'attribut à ajouter.
     * @param value Le valeur de l'attribut à ajouter.
     * @return self
     */
    public Builder objectAttribute(ObjectAttribute type, Object value) {
      ensureNotBuilt();
      objectAttributes.put(type, value);
      return this;
    }

    /**
     * @return L'objet de mise à jour construit avec ce constructeur.
     */
    public EntityDataUpdate build() {
      built = true;
      return new EntityDataUpdate(this);
    }

    private void ensureNotBuilt() {
      if (built) {
        throw new IllegalStateException("object has already been built");
      }
    }

  }

  private final int entityId;
  private final Map<EntityUpdateType, Object> data;
  private final Map<ObjectAttribute, Object> objectAttributes;

  private EntityDataUpdate(Builder builder) {
    entityId = builder.entityId;
    data = builder.data;
    objectAttributes = builder.objectAttributes;
  }

  /**
   * @return L'id de l'entité à mettre à jour.
   */
  public int getEntityId() {
    return entityId;
  }

  public MapPoint getPosition() {
    Object position = data.get(EntityUpdateType.POSITION);
    if (position == null) {
      throw new NoSuchElementException("position has not been set");
    }
    return (MapPoint) position;
  }

  public float getSpeedAngle() {
    Object speedAngle = data.get(EntityUpdateType.SPEED_ANGLE);
    if (speedAngle == null) {
      throw new NoSuchElementException("speedAngle has not been set");
    }
    return (float) speedAngle;
  }

  public float getSpeedNorm() {
    Object speedNorm = data.get(EntityUpdateType.SPEED_NORM);
    if (speedNorm == null) {
      throw new NoSuchElementException("speedNorm has not been set");
    }
    return (float) speedNorm;
  }

  public EntityType getEntityType() {
    Object entityType = data.get(EntityUpdateType.ENTITY_TYPE);
    if (entityType == null) {
      throw new NoSuchElementException("entityType has not been set");
    }
    return (EntityType) entityType;
  }

  public SpriteType getSpriteType() {
    Object sprite = data.get(EntityUpdateType.SPRITE_TYPE);
    if (sprite == null) {
      throw new NoSuchElementException("sprite has not been set");
    }
    return (SpriteType) sprite;
  }

  public Object getObjectAttribute(ObjectAttribute type) {
    return objectAttributes.get(type);
  }

  public Map<ObjectAttribute, Object> getObjectAttributes() {
    return Collections.unmodifiableMap(objectAttributes);
  }

  public boolean hasPosition() {
    return data.containsKey(EntityUpdateType.POSITION);
  }

  public boolean hasSpeedAngle() {
    return data.containsKey(EntityUpdateType.SPEED_ANGLE);
  }

  public boolean hasSpeedNorm() {
    return data.containsKey(EntityUpdateType.SPEED_NORM);
  }

  public boolean hasEntityType() {
    return data.containsKey(EntityUpdateType.ENTITY_TYPE);
  }

  public boolean hasSprite() {
    return data.containsKey(EntityUpdateType.SPRITE_TYPE);
  }

}
