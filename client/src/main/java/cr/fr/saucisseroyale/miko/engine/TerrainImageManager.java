package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import fr.delthas.uitest.Drawer;
import fr.delthas.uitest.Image;
import org.lwjgl.BufferUtils;

import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.function.BiFunction;

/**
 * Un gestionnaire des images correspondant au terrain de jeu, mis à jour à chaque changement de chunk.
 */
class TerrainImageManager {
  // TODO palette mode, old image modification, image deletion if too far away
  private BiFunction<Long, ChunkPoint, Chunk> chunkProvider;
  private ChunkData defaultChunk;
  private long tick = Long.MAX_VALUE;
  private Map<ChunkPoint, ChunkData> map = new HashMap<>(50);
  private Set<ChunkPoint> toUpdate = new HashSet<>();
  private ByteBuffer buffer = BufferUtils.createByteBuffer(256 * 256 * 3);

  public TerrainImageManager(BiFunction<Long, ChunkPoint, Chunk> chunkProvider) {
    this.chunkProvider = chunkProvider;
    defaultChunk = updateChunk(null, TerrainManager.getDefaultChunk());
  }

  public void updateTick(long newTick) {
    tick = newTick;
    if (toUpdate.isEmpty()) {
      return;
    }
    for (ChunkPoint chunkPoint : toUpdate) {
      Chunk chunk = chunkProvider.apply(tick, chunkPoint);
      updateChunk(chunkPoint, chunk);
    }
  }

  public void update(long tick, ChunkPoint chunkPoint) {
    ChunkData chunkData = map.get(chunkPoint);
    Chunk latestChunk = chunkProvider.apply(this.tick, chunkPoint);
    if (chunkData != null && latestChunk == chunkData.chunk) {
      return;
    }
    if (tick > this.tick) {
      toUpdate.add(chunkPoint);
      return;
    }
    updateChunk(chunkPoint, latestChunk);
  }

  public void drawChunk(Drawer drawer, ChunkPoint chunkPoint) {
    ChunkData chunkData = map.getOrDefault(chunkPoint, defaultChunk);
    if (chunkData.image == null) {
      drawer.setColor(chunkData.chunk.getDefaultType().getColor());
      drawer.fillRectangle(0, 0, 256, 256, false);
    } else {
      drawer.drawImage(0, 0, chunkData.image, false);
    }
  }

  private ChunkData updateChunk(ChunkPoint chunkPoint, Chunk chunk) {
    ChunkData chunkData;
    if (chunk.isUniform()) {
      chunkData = new ChunkData(chunk);
    } else {
      buffer.clear();
      for (int y = 255; y >= 0; y--) {
        for (int x = 0; x < 256; x++) {
          int color = chunk.getBlock(x, y).getColorInt();
          buffer.put((byte) (color >> 16 & 0xFF)).put((byte) (color >> 8 & 0xFF)).put((byte) (color & 0xFF));
        }
      }
      Image image = Image.createImageRaw(buffer, 256, 256, true);
      chunkData = new ChunkData(chunk, image);
    }
    if (chunkPoint != null) {
      map.put(chunkPoint, chunkData);
    }
    return chunkData;
  }

  private static class ChunkData {
    public final Chunk chunk;
    public final Image image;

    public ChunkData(Chunk chunk, Image image) {
      this.chunk = chunk;
      this.image = image;
    }

    public ChunkData(Chunk chunk) {
      this.chunk = chunk;
      image = null;
    }
  }
}
