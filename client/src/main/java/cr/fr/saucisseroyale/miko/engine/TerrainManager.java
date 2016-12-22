package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.TerrainType;

import java.util.Collections;
import java.util.function.BiConsumer;

/**
 * Un gestionnaire du terrain de jeu, stockant tous les chunks à tous les ticks, avec le principe de
 * {@link Snapshots}.
 *
 * @see Chunk
 * @see Snapshots
 */
class TerrainManager {
  private static final Chunk defaultChunk = new Chunk(TerrainType.UNKNOWN, Collections.emptyList());
  private SnapshotsMap<ChunkPoint, Chunk> map = new SnapshotsMap<>(5);
  private BiConsumer<Long, ChunkPoint> chunkUpdateConsumer;

  /**
   * Retourne le chunk que ce gestionnaire fournit si un chunk demandé n'a pas été défini.
   *
   * @return Le chunk par défaut.
   * @see #chunkDefined(ChunkPoint)
   */
  public static Chunk getDefaultChunk() {
    return defaultChunk;
  }

  /**
   * Ajoute un chunk à la carte de jeu, aux coordonnées spécifiées, au tick spécifié.
   *
   * @param tick     Le tick auquel ajouter le chunk.
   * @param position La position à laquelle ajouter le chunk.
   * @param chunk    Le chunk à ajouter.
   */
  public void setChunk(long tick, ChunkPoint position, Chunk chunk) {
    map.setSnapshot(tick, position, chunk);
    if (chunkUpdateConsumer != null) {
      chunkUpdateConsumer.accept(tick, position);
    }
  }

  /**
   * Retourne le chunk à l'endroit spécifié, au tick spécifié, ou null s'il n'existe pas de chunk à
   * ces coordonnées temporelles et spatiales.
   * <p>
   * Lorsque le chunk n'a pas été défini, renvoit un chunk par défaut ne comportant que des
   * {@link TerrainType#UNKNOWN}.
   *
   * @param tick     Le tick du chunk à renvoyer.
   * @param position La position du chunk à renvoyer.
   * @return Le chunk spécifié par les coordonnées et le tick, ou null s'il n'existe pas.
   */
  public Chunk getChunk(long tick, ChunkPoint position) {
    Chunk chunk = map.getSnapshot(tick, position);
    if (chunk == null) {
      return defaultChunk;
    }
    return chunk;
  }

  /**
   * Retourne vrai si le chunk à l'endroit spécifié a été défini.
   *
   * @param position La position où chercher si le chunk existe.
   * @return true si le chunk a été défini à cette position.
   * @see #getChunk(long, ChunkPoint)
   */
  public boolean chunkDefined(ChunkPoint position) {
    return map.getSnapshot(Long.MAX_VALUE, position) != null;
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

  /**
   * Définit le listener qui sera appelé à chaque mise à jour d'un chunk pour le notifier.
   *
   * @param chunkUpdateConsumer Le listener à définir.
   */
  public void setChunkUpdateConsumer(BiConsumer<Long, ChunkPoint> chunkUpdateConsumer) {
    this.chunkUpdateConsumer = chunkUpdateConsumer;
  }
}
