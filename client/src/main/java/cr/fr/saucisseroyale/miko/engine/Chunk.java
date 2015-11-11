package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.TerrainType;

import java.util.Arrays;
import java.util.List;

/**
 * Un bloc immutable de terrain de 256x256 cases.
 *
 * @see Block
 */
public final class Chunk {


  private final TerrainType defaultType;
  private final Block[] blocks;

  /**
   * @param defaultType Le type de case à utiliser par défaut lorsque la case n'a pas été ajoutée.
   * @param data Les cases à mettre dans le bloc.
   */
  public Chunk(TerrainType defaultType, List<Block> data) {
    this.defaultType = defaultType;
    // Utilisation d'un tableau pour la mise en cache
    // Un tableau et non pas 4 "buckets" pour avoir suffisament de cases à comparer (sinon le
    // cache est toujours presque vide et l'overhead des objets est trop grand)
    int size = data.size();
    blocks = new Block[size];
    // Copie pour maîtriser l'implémentation (et défensive)
    int i = 0;
    for (Block block : data) {
      blocks[i] = block;
      i++;
    }
  }

  /**
   * @param x La position en x de la case dans le bloc.
   * @param y La position en y de la case dans le bloc.
   * @return Le type de terrain en la case spécifiée.
   */
  public TerrainType getBlock(int x, int y) {
    if (x < 0 || x >= 1 << 8) {
      throw new IllegalArgumentException("x must be between 0 and 255 inclusive");
    }
    if (y < 0 || y >= 1 << 8) {
      throw new IllegalArgumentException("y must be between 0 and 255 inclusive");
    }
    for (Block block : blocks) {
      if (block.getX() == x && block.getY() == y) {
        return block.getType();
      }
    }
    return defaultType;
  }

  @Override
  public int hashCode() {
    final int prime = 31;
    int result = 1;
    result = prime * result + Arrays.hashCode(blocks);
    result = prime * result + (defaultType == null ? 0 : defaultType.hashCode());
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
    if (!(obj instanceof Chunk)) {
      return false;
    }
    Chunk other = (Chunk) obj;
    if (!Arrays.equals(blocks, other.blocks)) {
      return false;
    }
    if (defaultType != other.defaultType) {
      return false;
    }
    return true;
  }

}
