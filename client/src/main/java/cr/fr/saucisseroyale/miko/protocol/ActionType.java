package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type d'action.
 *
 */
public enum ActionType implements UniquelyIdentifiable {

  // @noformatting
  SEND_BALL(0, DataType.ONE_FLOAT);
  // @formatting

  static {
    IdSaver.register(ActionType.class);
  }

  private final int id;
  private final DataType dataType;

  private ActionType(int id, DataType dataType) {
    assert id < 1 << 16 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    assert dataType != null : "le type de données doivent être définis";
    this.id = id;
    this.dataType = dataType;
  }

  /**
   * @return L'identifiant correspondant au type d'action.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * @return Le type de données pour ce type d'action.
   */
  public DataType getDataType() {
    return dataType;
  }

  /**
   * Renvoit la valeur d'un code de type d'action.
   *
   * @param id L'identifiant correspondant au type de action.
   * @return Le type d'action, ou null s'il n'y a pas de type d'action correspondante.
   */
  public static ActionType getType(int id) {
    return IdSaver.getValue(ActionType.class, id);
  }
}
