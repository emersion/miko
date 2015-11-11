package cr.fr.saucisseroyale.miko.protocol;

/**
 * Un point immutable de la carte de jeu, indépendante de son type. Utilisé pour indiquer un
 * emplacement sans répéter (bx,by,x,y).
 *
 */
public final class MapPoint {

  private final int chunkX;
  private final int chunkY;
  private final int blockX;
  private final int blockY;

  /**
   * @param chunkX La coordonnée X du bloc.
   * @param chunkY La coordonnée Y du bloc.
   * @param blockX La coordonnée X de la case.
   * @param blockY La coordonnée Y de la case.
   */
  public MapPoint(int chunkX, int chunkY, int blockX, int blockY) {
    if (chunkX < -(1 << 15) || chunkX >= 1 << 15) {
      throw new IllegalArgumentException("chunkX must be between -32768 and 32767 inclusive");
    }
    if (chunkY < -(1 << 15) || chunkY >= 1 << 15) {
      throw new IllegalArgumentException("chunkY must be between -32768 and 32767 inclusive");
    }
    if (blockX < 0 || blockX >= 1 << 8) {
      throw new IllegalArgumentException("blockX must be between 0 and 255 inclusive");
    }
    if (blockY < 0 || blockY >= 1 << 8) {
      throw new IllegalArgumentException("blockY must be between 0 and 255 inclusive");
    }
    this.chunkX = chunkX;
    this.chunkY = chunkY;
    this.blockX = blockX;
    this.blockY = blockY;
  }

  /**
   * @return La coordonnée X du bloc.
   */
  public int getChunkX() {
    return chunkX;
  }

  /**
   * @return La coordonnée Y du bloc.
   */
  public int getChunkY() {
    return chunkY;
  }

  /**
   * @return La coordonnée X de la case.
   */
  public int getBlockX() {
    return blockX;
  }

  /**
   * @return La coordonnée Y de la case.
   */
  public int getBlockY() {
    return blockY;
  }

  @Override
  public int hashCode() {
    final int prime = 31;
    int result = 1;
    result = prime * result + blockX;
    result = prime * result + blockY;
    result = prime * result + chunkX;
    result = prime * result + chunkY;
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
    if (!(obj instanceof MapPoint)) {
      return false;
    }
    MapPoint other = (MapPoint) obj;
    if (blockX != other.blockX) {
      return false;
    }
    if (blockY != other.blockY) {
      return false;
    }
    if (chunkX != other.chunkX) {
      return false;
    }
    if (chunkY != other.chunkY) {
      return false;
    }
    return true;
  }

}
