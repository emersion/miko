package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de terminaison (envoyé lors de la terminaison de la connection).
 *
 */
public enum ExitType implements UniquelyIdentifiable {

  // @noformatting
  CLIENT_QUIT(0),
  SERVER_CLOSED(1),
  NETWORK_ERROR(2),
  PING_TIMEOUT(3),
  CLIENT_KICKED(4),
  CLIENT_BANNED(5);
  // @formatting

  static {
    IdSaver.register(ExitType.class);
  }

  private final int id;

  private ExitType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant au type de terminaison.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de terminaison.
   *
   * @param id L'identifiant correspondant au type de terminaison.
   * @return Le type de la terminaison, ou null s'il n'y a pas de type de terminaison correspondant.
   */
  public static ExitType getType(int id) {
    return IdSaver.getValue(ExitType.class, id);
  }
}
