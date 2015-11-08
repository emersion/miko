package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de données dans la mise à jour d'une entité
 *
 */
public enum EntityUpdateType implements UniquelyIdentifiable {


  // @noformatting
  POSITION(0),
  SPEED_ANGLE(1),
  SPEED_NORM(2),
  SPRITE(6),
  OBJECT_DATA(7);
  // @formatting

  static {
    IdSaver.register(EntityUpdateType.class);
  }

  private final int id;

  private EntityUpdateType(int id) {
    assert id < 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant au type de mise à jour.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de type de mise à jour.
   *
   * @param id L'identifiant correspondant au type de mise à jour.
   * @return Le type de mise à jour, ou null s'il n'y a pas de type de mise à jour correspondant.
   */
  public static EntityUpdateType getType(int id) {
    return IdSaver.getValue(EntityUpdateType.class, id);
  }
}
