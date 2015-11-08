package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.TerrainType;

/**
 * Une case immutable de terrain qui a une position et un type.
 *
 */
public final class Block {

  private final int x;
  private final int y;
  private final TerrainType type;

  /**
   * @param x La coordonnée X de la case à créer.
   * @param y La coordonnée Y de la case à créer.
   * @param type Le type de terrain de la case à créer.
   */
  public Block(int x, int y, TerrainType type) {
    if (x < 0 || x >= 1 << 8) {
      throw new IllegalArgumentException("x must be between 0 and 255 inclusive");
    }
    if (y < 0 || y >= 1 << 8) {
      throw new IllegalArgumentException("y must be between 0 and 255 inclusive");
    }
    this.x = x;
    this.y = y;
    this.type = type;
  }

  /**
   * @return La coordonnée X de la case.
   */
  public int getX() {
    return x;
  }

  /**
   * @return La coordonnée Y de la case.
   */
  public int getY() {
    return y;
  }

  /**
   * @return Le type de terrain de la case.
   */
  public TerrainType getType() {
    return type;
  }



}
