package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un sprite : une animation d'entité.
 * <p>
 * Objet immutable définissant un type d'animation d'une entité, sans posséder la logique nécessaire
 * au rendu de l'animation.
 *
 */
public enum Sprite implements UniquelyIdentifiable {

  // @noformatting
  ;
  // @formatting

  static {
    IdSaver.register(Sprite.class);
  }

  private final int id;

  private Sprite(int id) {
    assert id < 1 << 16 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant au sprite.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de sprite.
   *
   * @param id L'identifiant correspondant au sprite.
   * @return Le sprite, ou null s'il n'y a pas de sprite correspondant.
   */
  public static Sprite getType(int id) {
    return IdSaver.getValue(Sprite.class, id);
  }
}
