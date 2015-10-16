package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de réponse à une requête de création de compte.
 *
 */
public enum RegisterResponseType implements UniquelyIdentifiable {

  // @noformatting
  OK(0),
  USED_PSEUDO(1),
  INVALID_PSEUDO(2),
  INVALID_PASSWORD(3),
  TOO_MANY_TRIES(4),
  REGISTER_DISABLED(5);
  // @formatting

  static {
    IdSaver.register(RegisterResponseType.class);
  }

  private final int id;

  private RegisterResponseType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant à la réponse d'enregistrement.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de réponse d'enregistrement.
   *
   * @param id L'identifiant correspondant à la réponse d'enregistrement.
   * @return Le type de la réponse, ou null s'il n'y a pas de réponse correspondante.
   */
  public static RegisterResponseType getType(int id) {
    return IdSaver.getValue(RegisterResponseType.class, id);
  }
}
