package cr.fr.saucisseroyale.miko.protocol;


/**
 * Un bloc immutable de la carte de jeu, indépendant de son type. Utilisé pour indiquer un
 * emplacement sans répéter (bx,by).
 *
 */
public final class ChunkPoint {

  private final int chunkX;
  private final int chunkY;

  /**
   * @param chunkX La coordonnée X du bloc.
   * @param chunkY La coordonnée Y du bloc.
   */
  public ChunkPoint(int chunkX, int chunkY) {
    if (chunkX < -(1 << 15) || chunkX >= 1 << 15)
      throw new IllegalArgumentException("chunkX must be between -32768 and 32767 inclusive");
    if (chunkY < -(1 << 15) || chunkY >= 1 << 15)
      throw new IllegalArgumentException("chunkY must be between -32768 and 32767 inclusive");
    this.chunkX = chunkX;
    this.chunkY = chunkY;
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


}
