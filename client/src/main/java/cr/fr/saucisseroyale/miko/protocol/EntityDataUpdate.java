package cr.fr.saucisseroyale.miko.protocol;

import java.util.EnumMap;
import java.util.EnumSet;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Set;

/**
 * Un objet stockant des valeurs pour mettre à jour une entité.
 * 
 */
public final class EntityDataUpdate {

  private final int entityId;
  private final Map<EntityUpdateType, Object> data = new EnumMap<>(EntityUpdateType.class);

  /**
   * @param entityId L'id de l'entité à mettre à jour.
   */
  public EntityDataUpdate(int entityId) {
    this.entityId = entityId;
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

  /**
   * 
   * @param position La nouvelle position de l'entité.
   */
  public void setPosition(MapPoint position) {
    data.put(EntityUpdateType.POSITION, position);
  }

  /**
   * 
   * @param speedAngle Le nouvel angle du vecteur vitesse.
   */
  public void setSpeedAngle(float speedAngle) {
    data.put(EntityUpdateType.SPEED_ANGLE, speedAngle);
  }

  /**
   * 
   * @param speedNorm La nouvelle norme du vecteur vitesse.
   */
  public void setSpeedNorm(float speedNorm) {
    data.put(EntityUpdateType.SPEED_NORM, speedNorm);
  }

  /**
   * 
   * @param attributes Les nouveaux attributs de l'objet.
   */
  public void setObjectAttributes(List<ObjectAttribute> attributes) {
    data.put(EntityUpdateType.OBJECT_DATA, EnumSet.copyOf(attributes));
  }

}
