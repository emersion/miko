package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type d'entité.
 *
 */
public enum EntityType implements UniquelyIdentifiable {

  // @noformatting
  PLAYER(0),
  BALL(1);
  // @formatting

  static {
    IdSaver.register(EntityType.class);
  }

  private final int id;

  private EntityType(int id) {
    assert id < 1 << 16 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant au type d'entité.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de type d'entité.
   *
   * @param id L'identifiant correspondant au type d'entité.
   * @return Le type d'entité, ou null s'il n'y a pas de type d'entité correspondant.
   */
  public static EntityType getType(int id) {
    return IdSaver.getValue(EntityType.class, id);
  }
}
