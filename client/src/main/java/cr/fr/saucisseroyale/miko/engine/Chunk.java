package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.TerrainType;

import java.util.ArrayList;
import java.util.List;

/**
 * Un bloc de terrain de 256x256 cases.
 * 
 * @see Block
 */
public final class Chunk {


  private final TerrainType defaultType;
  private final List<Block> blocks;

  /**
   * @param defaultType Le type de case à utiliser par défaut lorsque la case n'a pas été ajoutée.
   * @param data Les cases à mettre dans le bloc.
   */
  public Chunk(TerrainType defaultType, List<Block> data) {
    this.defaultType = defaultType;
    // ArrayList pour mise en cache du tableau sous-jacent
    // Une LinkedList ou une Hashmap n'en auraient pas profité
    // Une ArrayList et non pas 4 "buckets" pour avoir suffisament de cases à comparer (sinon le
    // cache est toujours presque vide et l'overhead des objets est trop grand)
    // 5 cases supplémentaires en cas de petit changement pour ne pas tout redimensionner
    blocks = new ArrayList<>(data.size() + 5);
    // Copie pour maîtriser l'implémentation (et défensive)
    for (Block block : data) {
      blocks.add(block);
    }
  }

  /**
   * @param x La position en x de la case dans le bloc.
   * @param y La position en y de la case dans le bloc.
   * @param type Le type de terrain en la case spécifiée.
   */
  public void putBlock(int x, int y, TerrainType type) {
    putBlock(new Block(x, y, type));
  }

  /**
   * @param block La case à ajouter au bloc de cases.
   */
  public void putBlock(Block block) {
    if (block.getType() == defaultType)
      return;
    if (block.getX() < 0 || block.getX() >= 1 << 8)
      throw new IllegalArgumentException("x must be between 0 and 255 inclusive");
    if (block.getY() < 0 || block.getY() >= 1 << 8)
      throw new IllegalArgumentException("y must be between 0 and 255 inclusive");
    for (int i = 0, n = blocks.size(); i < n; i++) {
      Block oldBlock = blocks.get(i);
      if (oldBlock.getX() == block.getX() && oldBlock.getY() == block.getY()) {
        blocks.set(i, block);
        return;
      }
    }
    blocks.add(block);
  }


  /**
   * @param x La position en x de la case dans le bloc.
   * @param y La position en y de la case dans le bloc.
   * @return Le type de terrain en la case spécifiée.
   */
  public TerrainType getBlock(int x, int y) {
    if (x < 0 || x >= 1 << 8)
      throw new IllegalArgumentException("x must be between 0 and 255 inclusive");
    if (y < 0 || y >= 1 << 8)
      throw new IllegalArgumentException("y must be between 0 and 255 inclusive");
    for (Block block : blocks) {
      if (block.getX() == x && block.getY() == y)
        return block.getType();
    }
    return defaultType;
  }

  // TODO DEBUG
  public void dumpSize() {
    System.out.println(blocks.size());
  }
}
