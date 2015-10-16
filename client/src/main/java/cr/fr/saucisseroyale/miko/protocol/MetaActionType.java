/**
 * 
 */
package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de méta-action (action liée à un joueur et pas au jeu).
 *
 */
public enum MetaActionType implements UniquelyIdentifiable {


  // @noformatting
  PLAYER_JOINED(0),
  PLAYER_LEFT(1);
  // @formatting

  static {
    IdSaver.register(MetaActionType.class);
  }

  private final int id;

  private MetaActionType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant au type de méta-action.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de méta-action.
   *
   * @param id L'identifiant correspondant au type de méta-action.
   * @return Le type de la méta-action, ou null s'il n'y a pas de type de méta-action correspondant.
   */
  public static MetaActionType getType(int id) {
    return IdSaver.getValue(MetaActionType.class, id);
  }
}
