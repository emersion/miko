package cr.fr.saucisseroyale.miko.protocol;

import java.util.EnumMap;
import java.util.EnumSet;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Set;

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
    private boolean built = false;

    /**
     * @param entityId L'id de l'entité à mettre à jour.
     */
    public Builder(int entityId) {
      if (entityId < 0 || entityId >= 1 << 16)
        throw new IllegalArgumentException("entityId must be between 0 and 65535 inclusive");
      this.entityId = entityId;
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
     * @param attributes Les nouveaux attributs de l'objet.
     * @return self
     */
    public Builder objectAttributes(List<ObjectAttribute> attributes) {
      ensureNotBuilt();
      data.put(EntityUpdateType.OBJECT_DATA, EnumSet.copyOf(attributes));
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
      if (built)
        throw new IllegalStateException("object has already been built");
    }

  }

  private final int entityId;
  private final Map<EntityUpdateType, Object> data;

  private EntityDataUpdate(Builder builder) {
    entityId = builder.entityId;
    data = builder.data;
  }

  /**
   * @return L'id de l'entité à mettre à jour.
   */
  public int getEntityId() {
    return entityId;
  }

  public MapPoint getPosition() {
    Object position = data.get(EntityUpdateType.POSITION);
    if (position == null)
      throw new NoSuchElementException("position has not been set");
    return (MapPoint) position;
  }

  public float getSpeedAngle() {
    Object speedAngle = data.get(EntityUpdateType.SPEED_ANGLE);
    if (speedAngle == null)
      throw new NoSuchElementException("speedAngle has not been set");
    return (float) speedAngle;
  }

  public float getSpeedNorm() {
    Object speedNorm = data.get(EntityUpdateType.SPEED_NORM);
    if (speedNorm == null)
      throw new NoSuchElementException("speedNorm has not been set");
    return (float) speedNorm;
  }

  public Set<ObjectAttribute> getObjectAttributes() {
    Object attributesRaw = data.get(EntityUpdateType.OBJECT_DATA);
    if (attributesRaw == null)
      throw new NoSuchElementException("objectAttributes has not been set");
    // on sait que c'est du bon type
    @SuppressWarnings("unchecked")
    Set<ObjectAttribute> attributes = (Set<ObjectAttribute>) attributesRaw;
    return attributes;
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

  public boolean hasObjectAttributes() {
    return data.containsKey(EntityUpdateType.OBJECT_DATA);
  }

}
