package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;

import java.awt.*;

/**
 * Un type de case de terrain.
 */
public enum TerrainType implements UniquelyIdentifiable {


  // @noformatting
  WHITE_GROUND(0, 255, 255, 255),
  BLACK_WALL(1, 0, 0, 0),
  UNKNOWN(255, 127, 127, 127);
  // @formatting

  static {
    IdSaver.register(TerrainType.class);
  }

  private final int id;
  private final int colorInt;
  private final Color color;

  TerrainType(int id, int red, int green, int blue) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
    colorInt = red << 16 | green << 8 | blue;
    color = new Color(colorInt, false);
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

  /**
   * @return L'identifiant correspondant au type de case.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * @return L'entier représentant la couleur de cette case sous la forme 0x00RRGGBB.
   */
  public int getColorInt() {
    return colorInt;
  }

  /**
   * @return La couleur de cette case.
   */
  public Color getColor() {
    return color;
  }
}
