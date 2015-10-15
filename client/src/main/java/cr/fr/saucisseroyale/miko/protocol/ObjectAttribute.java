/**
 * 
 */
package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Une caractéristique d'un objet lors d'une mise à jour.
 *
 */
public enum ObjectAttribute implements UniquelyIdentifiable {

  // @noformatting
  DISABLED(0),
  ENABLED(1);
  // @formatting

  static {
    IdSaver.register(ObjectAttribute.class, 1 << 8);
  }

  private final int id;

  private ObjectAttribute(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant à la caractéristique.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un identifiant de caractéristique.
   *
   * @param id L'identifiant correspondant à la caractéristique.
   * @return La caractéristique, ou null s'il n'y a pas de caractéristique correspondante.
   */
  public static ObjectAttribute getType(int id) {
    return IdSaver.getValue(ObjectAttribute.class, id);
  }
}
