package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.engine.BlockPoint;
import cr.fr.saucisseroyale.miko.util.MikoMath;

/**
 * Un point immutable de la carte de jeu, indépendante de son type. Utilisé pour indiquer un
 * emplacement sans répéter (bx,by,x,y).
 */
public final class TerrainPoint {
  private final int x;
  private final int y;

  /**
   * @param chunkX La coordonnée X du bloc.
   * @param chunkY La coordonnée Y du bloc.
   * @param blockX La coordonnée X de la case dans le bloc.
   * @param blockY La coordonnée Y de la case dans le bloc.
   */
  public TerrainPoint(int chunkX, int chunkY, int blockX, int blockY) {
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
    x = chunkX * 256 + blockX;
    y = chunkY * 256 + blockY;
  }

  /**
   * @param x La coordonnée X du point.
   * @param y La coordonnée Y du point.
   */
  public TerrainPoint(int x, int y) {
    this.x = x;
    this.y = y;
  }

  /**
   * @return La coordonnée X du point.
   */
  public int getX() {
    return x;
  }

  /**
   * @return La coordonnée Y du point.
   */
  public int getY() {
    return y;
  }

  /**
   * @return La coordonnée X du bloc.
   */
  public int getChunkX() {
    return MikoMath.quotient(x, 256);
  }

  /**
   * @return La coordonnée Y du bloc.
   */
  public int getChunkY() {
    return MikoMath.quotient(y, 256);
  }

  /**
   * @return Un ChunkPoint correspondant au bloc de ce point.
   */
  public ChunkPoint getChunkPoint() {
    return new ChunkPoint(getChunkX(), getChunkY());
  }

  /**
   * @return La coordonnée X de la case dans le bloc.
   */
  public int getBlockX() {
    return MikoMath.modulo(x, 256);
  }

  /**
   * @return La coordonnée Y de la case dans le bloc.
   */
  public int getBlockY() {
    return MikoMath.modulo(y, 256);
  }

  /**
   * @return Un BlockPoint correspondant à la case de ce point dans le bloc.
   */
  public BlockPoint getBlockPoint() {
    return new BlockPoint(getBlockX(), getBlockY());
  }

  /**
   * Retourne ce point après translation de (deltaX;deltaY) (déplacé de deltaX en abscisse et de
   * deltaY en ordonnée).
   *
   * @param deltaX La translation en X à effectuer.
   * @param deltaY La translation en Y à effectuer.
   * @return Le point après translation.
   */
  public TerrainPoint getTranslated(int deltaX, int deltaY) {
    return new TerrainPoint(x + deltaX, y + deltaY);
  }

  @Override
  public int hashCode() {
    final int prime = 31;
    int result = 1;
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
    if (!(obj instanceof TerrainPoint)) {
      return false;
    }
    TerrainPoint other = (TerrainPoint) obj;
    if (x != other.x) {
      return false;
    }
    return y == other.y;
  }
}
