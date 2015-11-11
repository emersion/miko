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

  @Override
  public int hashCode() {
    final int prime = 31;
    int result = 1;
    result = prime * result + (type == null ? 0 : type.hashCode());
    result = prime * result + x;
    result = prime * result + y;
    return result;
  }

  @Override
  public boolean equals(Object obj) {
    if (this == obj) {
      return true;
    }
    if (obj == null) {
      return false;
    }
    if (!(obj instanceof Block)) {
      return false;
    }
    Block other = (Block) obj;
    if (type != other.type) {
      return false;
    }
    if (x != other.x) {
      return false;
    }
    if (y != other.y) {
      return false;
    }
    return true;
  }

}
