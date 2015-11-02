package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

/**
 * Un type de case de terrain.
 *
 */
public enum TerrainType implements UniquelyIdentifiable {


  // @noformatting
  // TODO
  NORMAL(0),
  ANORMAL(1);
  // @formatting

  static {
    IdSaver.register(TerrainType.class);
  }

  private final int id;

  private TerrainType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
   * @return L'identifiant correspondant au type de case.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit la valeur du type de case.
   *
   * @param id L'identifiant correspondant au type de case.
   * @return Le type de la case, ou null s'il n'y a pas de type de case correspondant.
   */
  public static TerrainType getType(int id) {
    return IdSaver.getValue(TerrainType.class, id);
  }
}
