package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un attribut d'un objet lors d'une mise à jour d'attributs.
 */
public enum ObjectAttribute implements UniquelyIdentifiable {

  // @noformatting
  TICKS_LEFT(0, DataType.ONE_SHORT),
  HEALTH(1, DataType.ONE_SHORT),
  SENDER(2, DataType.ONE_ENTITY),
  COOLDOWN_ONE(30000, DataType.ONE_SHORT);
  // @formatting

  static {
    IdSaver.register(ObjectAttribute.class);
  }

  private final int id;
  private final DataType dataType;

  ObjectAttribute(int id, DataType dataType) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    assert dataType != null : "le type de données doivent être définis";
    this.id = id;
    this.dataType = dataType;
  }

  /**
   * Renvoit la valeur d'un identifiant d'attribut.
   *
   * @param id L'identifiant correspondant à l'attribut.
   * @return L'attribut, ou null s'il n'y a pas d'attribut correspondante.
   */
  public static ObjectAttribute getType(int id) {
    return IdSaver.getValue(ObjectAttribute.class, id);
  }

  /**
   * @return L'identifiant correspondant à l'attribut.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * @return Le type de données de l'attribut.
   */
  public DataType getDataType() {
    return dataType;
  }
}
