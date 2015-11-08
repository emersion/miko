package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.ListIterator;
import java.util.Map;
import java.util.Set;

/**
 * Un gestionnaire du terrain de jeu, stockant tous les chunks à tous les ticks.
 * 
 * @see Chunk
 */
class TerrainMap {

  private static final int chunkListCapacity = 5;
  private Map<ChunkPoint, List<Pair<Chunk, Long>>> map;

  /**
   * Ajoute un chunk à la carte de jeu, aux coordonnées spécifiées, au tick spécifié.
   * <p>
   * Le chunk sera aussi ajouté à tous les ticks suivants celui indiqué, jusqu'à ce qu'un autre
   * chunk soit ajouté. Autrement dit, les ticks ne possédant pas de chunks sont remplis avec les
   * chunks les plus récents avant ce tick.
   * 
   * @param tick Le tick auquel ajouter le chunk.
   * @param position La position à laquelle ajouter le chunk.
   * @param chunk Le chunk à ajouter.
   */
  public void addChunk(long tick, ChunkPoint position, Chunk chunk) {
    List<Pair<Chunk, Long>> pairList = map.get(position);
    if (pairList == null) {
      pairList = new ArrayList<>(chunkListCapacity);
      map.put(position, pairList);
    }
    ListIterator<Pair<Chunk, Long>> pairListIterator = pairList.listIterator();
    while (pairListIterator.hasNext()) {
      Pair<Chunk, Long> pair = pairListIterator.next();
      if (pair.getSecond() >= tick) {
        if (pair.getSecond() == tick) {
          pairListIterator.remove();
        }
        pairListIterator.previous();
        pairListIterator.add(new Pair<>(chunk, tick));
        return;
      }
    }
    pairList.add(new Pair<>(chunk, tick));
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
    List<Pair<Chunk, Long>> pairList = map.get(position);
    if (pairList == null) {
      return null;
    }
    ListIterator<Pair<Chunk, Long>> pairListIterator = pairList.listIterator();
    while (pairListIterator.hasNext()) {
      Pair<Chunk, Long> pair = pairListIterator.next();
      if (pair.getSecond() == tick) {
        return pair.getFirst();
      }
      if (pair.getSecond() > tick) {
        pairListIterator.previous();
        if (!pairListIterator.hasPrevious()) {
          return null;
        }
        return pairListIterator.previous().getFirst();
      }
    }
    if (pairList.isEmpty()) {
      return null;
    }
    return pairList.get(pairList.size() - 1).getFirst();
  }

  /**
   * Indique au gestionnaire de terrain que les chunks appartenant à des ticks avant ou égaux au
   * tick spécifié ne seront plus jamais demandés et peuvent être supprimés.
   * <p>
   * Cette méthode n'est qu'une indication pour le gestionnaire, qui peut décider de disposer les
   * ressources associées aux ticks, ou non.
   * 
   * @param tick Le tick (inclus) jusqu'auquel les chunks ne seront plus demandés.
   */
  public void disposeTo(long tick) {
    // NOTE on pourrait ignorer ce call la plupart du temps pour éviter de parcourir toute la map à
    // chaque fois ; il faudrait sans doute modifier la capacité standard
    Set<Map.Entry<ChunkPoint, List<Pair<Chunk, Long>>>> entrySet = map.entrySet();
    Iterator<Map.Entry<ChunkPoint, List<Pair<Chunk, Long>>>> entryIterator = entrySet.iterator();
    while (entryIterator.hasNext()) {
      Map.Entry<ChunkPoint, List<Pair<Chunk, Long>>> entry = entryIterator.next();
      List<Pair<Chunk, Long>> pairList = entry.getValue();
      Iterator<Pair<Chunk, Long>> pairListIterator = pairList.listIterator();
      while (pairListIterator.hasNext()) {
        Pair<Chunk, Long> pair = pairListIterator.next();
        if (pair.getSecond() > tick) {
          break;
        }
        pairListIterator.remove();
      }
      if (pairList.isEmpty()) {
        entryIterator.remove();
      }
    }
  }
}
