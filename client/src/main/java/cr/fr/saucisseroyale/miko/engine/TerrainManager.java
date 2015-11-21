package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;

/**
 * Un gestionnaire du terrain de jeu, stockant tous les chunks à tous les ticks, avec le principe de
 * {@link Snapshots}.
 *
 * @see Chunk
 * @see Snapshots
 */
class TerrainManager {

  private SnapshotsMap<ChunkPoint, Chunk> map = new SnapshotsMap<>(5);

  /**
   * Ajoute un chunk à la carte de jeu, aux coordonnées spécifiées, au tick spécifié.
   *
   * @param tick Le tick auquel ajouter le chunk.
   * @param position La position à laquelle ajouter le chunk.
   * @param chunk Le chunk à ajouter.
   */
  public void setChunk(long tick, ChunkPoint position, Chunk chunk) {
    map.setSnapshot(tick, position, chunk);
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
    return map.getSnapshot(tick, position);
  }

  /**
   * Indique au gestionnaire de terrain que les chunks appartenant à des ticks avant ou égaux au
   * tick spécifié ne seront plus jamais demandés et peuvent être supprimés.
   *
   * @param tick Le tick (inclus) jusqu'auquel les chunks ne seront plus demandés.
   */
  public void disposeUntilTick(long tick) {
    map.disposeUntilTick(tick);
  }
}
