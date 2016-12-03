package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.engine.Hitbox;
import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un sprite : une animation d'entité.
 * <p>
 * Objet immutable définissant un type d'animation d'une entité, sans posséder la logique nécessaire
 * au rendu de l'animation.
 */
public enum SpriteType implements UniquelyIdentifiable {

  // @noformatting
  PLACEHOLDER(0, Hitbox.newNullHitbox()),
  PLAYER(1, Hitbox.newCircleHitbox(10)),
  BALL(2, Hitbox.newCircleHitbox(10));
  // @formatting

  static {
    IdSaver.register(SpriteType.class);
  }

  private final int id;
  private final Hitbox hitbox;

  SpriteType(int id, Hitbox hitbox) {
    assert id < 1 << 16 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    assert hitbox != null : "la hitbox doit être définie";
    this.id = id;
    this.hitbox = hitbox;
  }

  /**
   * Renvoit la valeur d'un code de sprite.
   *
   * @param id L'identifiant correspondant au sprite.
   * @return Le sprite, ou null s'il n'y a pas de sprite correspondant.
   */
  public static SpriteType getType(int id) {
    return IdSaver.getValue(SpriteType.class, id);
  }

  /**
   * @return L'identifiant correspondant au sprite.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * @return La hitbox correspondant au sprite.
   */
  public Hitbox getHitbox() {
    return hitbox;
  }
}
