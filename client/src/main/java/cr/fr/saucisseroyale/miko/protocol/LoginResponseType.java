package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de réponse à une requête de connection.
 */
public enum LoginResponseType implements UniquelyIdentifiable {

  // @noformatting
  OK(0),
  UNKNOWN_PSEUDO(1),
  WRONG_PASSWORD(2),
  TOO_MANY_TRIES(3),
  ALREADY_CONNECTED(4),
  PLAYER_LIMIT_REACHED(5);
  // @formatting

  static {
    IdSaver.register(LoginResponseType.class);
  }

  private final int id;

  LoginResponseType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * Renvoit la valeur d'un code de réponse de connection.
   *
   * @param id L'identifiant correspondant à la réponse de connection.
   * @return Le type de la réponse, ou null s'il n'y a pas de réponse correspondante.
   */
  public static LoginResponseType getType(int id) {
    return IdSaver.getValue(LoginResponseType.class, id);
  }

  /**
   * @return L'identifiant correspondant à la réponse de connection.
   */
  @Override
  public int getId() {
    return id;
  }
}
