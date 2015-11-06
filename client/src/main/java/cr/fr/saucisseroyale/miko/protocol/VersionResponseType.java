package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de réponse à un envoi de version.
 *
 */
public enum VersionResponseType implements UniquelyIdentifiable {

  // @noformatting
  OK(0),
  CLIENT_OUTDATED(1),
  SERVER_OUTDATED(2);
  // @formatting

  static {
    IdSaver.register(VersionResponseType.class);
  }

  private final int id;

  private VersionResponseType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant à la réponse de version.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur d'un code de réponse de version.
   *
   * @param id L'identifiant correspondant à la réponse de version.
   * @return Le type de la réponse, ou null s'il n'y a pas de réponse correspondante.
   */
  public static VersionResponseType getType(int id) {
    return IdSaver.getValue(VersionResponseType.class, id);
  }
}
