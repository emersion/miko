package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.NoSuchElementException;

/**
 * Une action immutable, ayant un type et des paramètres.
 */
public final class Action {
  private final ActionType type;
  private final Object parameter;

  public Action(ActionType type, Object value) {
    this.type = type;
    switch (type.getDataType()) {
      case VOID:
        if (value != null) {
          throw newConstructorValueException();
        }
        parameter = null;
        break;
      case ONE_FLOAT:
        if (!(value instanceof Float)) {
          throw newConstructorValueException();
        }
        parameter = value;
        break;
      case ONE_SHORT:
      case ONE_ENTITY:
        if (!(value instanceof Integer)) {
          throw newConstructorValueException();
        }
        int intValue = (int) value;
        if (type.getDataType() == DataType.ONE_ENTITY && intValue < 0 || intValue >= 1 << 16) {
          throw new IllegalArgumentException("entityId must be between 0 and 65535 inclusive");
        }
        parameter = value;
        break;
      case ONE_TERRAIN:
        if (!(value instanceof TerrainPoint)) {
          throw newConstructorValueException();
        }
        parameter = value;
        break;
      case PAIR_FLOAT_ENTITY:
        if (!(value instanceof Pair)) {
          throw newConstructorValueException();
        }
        @SuppressWarnings("rawtypes")
        Pair pair = (Pair) value;
        if (!(pair.getFirst() instanceof Float) || !(pair.getSecond() instanceof Integer)) {
          throw newConstructorValueException();
        }
        parameter = value;
        break;
      default:
        throw new IllegalArgumentException("unknown data type: " + type);
    }
  }

  private static RuntimeException newConstructorValueException() {
    return new IllegalArgumentException("wrong parameter type");
  }

  private static RuntimeException newGetFieldException() {
    return new NoSuchElementException("this object has no such field");
  }

  /**
   * @return Le type de cette action.
   */
  public ActionType getType() {
    return type;
  }

  /**
   * @return Le type de données de cette action.
   */
  public DataType getDataType() {
    return type.getDataType();
  }

  /**
   * @return La valeur float de cette action.
   */
  public float getFloatValue() {
    if (type.getDataType() != DataType.ONE_FLOAT) {
      throw newGetFieldException();
    }
    return (float) parameter;
  }

  /**
   * @return La valeur entityId de cette action.
   */
  public int getEntityIdValue() {
    if (type.getDataType() != DataType.ONE_ENTITY) {
      throw newGetFieldException();
    }
    return (int) parameter;
  }

  public Pair<Float, Integer> getPairFloatEntityIdValue() {
    if (type.getDataType() != DataType.PAIR_FLOAT_ENTITY) {
      throw newGetFieldException();
    }
    @SuppressWarnings("unchecked")
    Pair<Float, Integer> pair = (Pair<Float, Integer>) parameter;
    return pair;
  }

  /**
   * @return La valeur short de cette action.
   */
  public int getShortValue() {
    if (type.getDataType() != DataType.ONE_SHORT) {
      throw newGetFieldException();
    }
    return (int) parameter;
  }

  /**
   * @return La valeur TerrainPoint de cette action.
   */
  public TerrainPoint getTerrainPointValue() {
    if (type.getDataType() != DataType.ONE_TERRAIN) {
      throw newGetFieldException();
    }
    return (TerrainPoint) parameter;
  }
}
