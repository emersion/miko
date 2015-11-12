package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;

import java.util.Iterator;
import java.util.Map;
import java.util.Set;

/**
 * Un gestionnaire du terrain de jeu, stockant tous les chunks à tous les ticks, avec le principe de
 * {@link Snapshots}.
 *
 * @see Chunk
 * @see Snapshots
 */
class TerrainManager {

  private static final int chunkListCapacity = 5;
  private Map<ChunkPoint, Snapshots<Chunk>> map;

  /**
   * Ajoute un chunk à la carte de jeu, aux coordonnées spécifiées, au tick spécifié.
   *
   * @param tick Le tick auquel ajouter le chunk.
   * @param position La position à laquelle ajouter le chunk.
   * @param chunk Le chunk à ajouter.
   */
  public void addChunk(long tick, ChunkPoint position, Chunk chunk) {
    Snapshots<Chunk> snapshots = map.get(position);
    if (snapshots == null) {
      snapshots = new Snapshots<>(chunkListCapacity);
      map.put(position, snapshots);
    }
    snapshots.setSnapshot(tick, chunk);
  }

  /**
   * Retourne le chunk à l'endroit spécifié, au tick spécifié, ou null s'il n'existe pas de chunk à
   * ces coordonnées temporelles et spatiales.
   *
   * @param tick Le tick du chunk à renvoyer.
   * @param position La position du chunk à renvoyer.
   * @return Le chunk spécifié par les coordonnées et le tick, ou null s'il n'existe pas.
   */
  public Chunk getChunk(long tick, ChunkPoint position) {
    Snapshots<Chunk> snapshots = map.get(position);
    if (snapshots == null) {
      return null;
    }
    return snapshots.getSnapshot(tick);
  }

  /**
   * Indique au gestionnaire de terrain que les chunks appartenant à des ticks avant ou égaux au
   * tick spécifié ne seront plus jamais demandés et peuvent être supprimés.
   *
   * @param tick Le tick (inclus) jusqu'auquel les chunks ne seront plus demandés.
   */
  public void disposeUntilTick(long tick) {
    // NOTE on pourrait ignorer ce call la plupart du temps pour éviter de parcourir toute la map à
    // chaque fois ; il faudrait sans doute modifier la capacité standard
    Set<Map.Entry<ChunkPoint, Snapshots<Chunk>>> entrySet = map.entrySet();
    Iterator<Map.Entry<ChunkPoint, Snapshots<Chunk>>> entryIterator = entrySet.iterator();
    while (entryIterator.hasNext()) {
      Map.Entry<ChunkPoint, Snapshots<Chunk>> entry = entryIterator.next();
      Snapshots<Chunk> snapshots = entry.getValue();
      snapshots.disposeUntilTick(tick);
      if (snapshots.isEmpty()) {
        entryIterator.remove();
      }
    }
  }
}
