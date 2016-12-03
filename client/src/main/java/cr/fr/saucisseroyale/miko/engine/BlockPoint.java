package cr.fr.saucisseroyale.miko.engine;

/**
 * Un point immutable dans un chunk, <b>exprimé en float</b>. Utilisé pour indiquer un emplacement
 * sans répéter (x,y).
 */
public final class BlockPoint {
  private final float blockX;
  private final float blockY;

  /**
   * @param blockX La coordonnée X de la case.
   * @param blockY La coordonnée Y de la case.
   */
  public BlockPoint(float blockX, float blockY) {
    if (blockX < 0 || blockX >= 1 << 8) {
      throw new IllegalArgumentException("blockX must be between 0 and 255 inclusive");
    }
    if (blockY < 0 || blockY >= 1 << 8) {
      throw new IllegalArgumentException("blockY must be between 0 and 255 inclusive");
    }
    this.blockX = blockX;
    this.blockY = blockY;
  }

  /**
   * @return La coordonnée X de la case.
   */
  public float getBlockX() {
    return blockX;
  }

  /**
   * @return La coordonnée Y de la case.
   */
  public float getBlockY() {
    return blockY;
  }

  @Override
  public int hashCode() {
    final int prime = 31;
    int result = 1;
    result = prime * result + Float.floatToIntBits(blockX);
    result = prime * result + Float.floatToIntBits(blockY);
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
    if (!(obj instanceof BlockPoint)) {
      return false;
    }
    BlockPoint other = (BlockPoint) obj;
    if (Float.floatToIntBits(blockX) != Float.floatToIntBits(other.blockX)) {
      return false;
    }
    return Float.floatToIntBits(blockY) == Float.floatToIntBits(other.blockY);
  }
}
