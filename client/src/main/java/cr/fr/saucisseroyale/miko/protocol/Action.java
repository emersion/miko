package cr.fr.saucisseroyale.miko.protocol;

import java.util.NoSuchElementException;

/**
 * Une action immutable, ayant un type et des paramètres.
 *
 */
public final class Action {

  private final ActionType type;
  private final Object parameter;

  public Action(ActionType type) {
    if (type.getParametersType() != ActionParametersType.VOID) {
      throw newConstructorValueException();
    }
    this.type = type;
    parameter = null;
  }

  public Action(ActionType type, float value) {
    if (type.getParametersType() != ActionParametersType.FLOAT) {
      throw newConstructorValueException();
    }
    this.type = type;
    parameter = value;
  }

  public Action(ActionType type, int entityId) {
    if (type.getParametersType() != ActionParametersType.ENTITY_ID) {
      throw newConstructorValueException();
    }
    if (entityId < 0 || entityId >= 1 << 16)
      throw new IllegalArgumentException("entityId must be between 0 and 65535 inclusive");
    this.type = type;
    parameter = entityId;
  }

  public Action(ActionType type, MapPoint mapPoint) {
    if (type.getParametersType() != ActionParametersType.MAP_POINT) {
      throw newConstructorValueException();
    }
    this.type = type;
    parameter = mapPoint;
  }

  /**
   * @return The type of this action.
   */
  public ActionType getType() {
    return type;
  }

  /**
   * @return The type of the parameters of this action.
   */
  public ActionParametersType getParameterType() {
    return type.getParametersType();
  }

  /**
   * @return The float value stored in this action.
   */
  public float getFloatValue() {
    if (type.getParametersType() != ActionParametersType.FLOAT) {
      throw newGetFieldException();
    }
    return (float) parameter;
  }

  /**
   * @return The entityId value stored in this action.
   */
  public int getEntityIdValue() {
    if (type.getParametersType() != ActionParametersType.ENTITY_ID) {
      throw newGetFieldException();
    }
    return (int) parameter;
  }

  /**
   * @return The map point value stored in this action.
   */
  public MapPoint getMapPointValue() {
    if (type.getParametersType() != ActionParametersType.MAP_POINT) {
      throw newGetFieldException();
    }
    return (MapPoint) parameter;
  }

  private static final RuntimeException newConstructorValueException() {
    return new IllegalArgumentException("wrong parameter type");
  }

  private static final RuntimeException newGetFieldException() {
    return new NoSuchElementException("this object has no such field");
  }

}
