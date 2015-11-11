package cr.fr.saucisseroyale.miko.engine;

import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.stream.IntStream;

/**
 * Un gestionnaire de joueurs, stockant tous les joueurs à tous les ticks, avec le principe de
 * {@link Snapshots}.
 *
 * @see Snapshots
 */
class PlayerManager {

  private static final int chunkListCapacity = 5;
  private Map<Integer, Snapshots<String>> map;

  /**
   * Ajoute un joueur à la liste de joueurs, au tick spécifié.
   *
   * @param tick Le tick auquel ajouter le joueur.
   * @param entityId L'entityId du joueur à ajouter.
   * @param name Le nom du joueur à ajouter.
   */
  public void addPlayer(long tick, int entityId, String name) {
    Snapshots<String> snapshots = map.get(entityId);
    if (snapshots == null) {
      snapshots = new Snapshots<>(chunkListCapacity);
      map.put(entityId, snapshots);
    }
    snapshots.setSnapshot(tick, name);
  }

  /**
   * Retourne le joueur au tick spécifié, ou null s'il n'existe pas de joueur à ce tick.
   *
   * @param tick Le tick du joueur à renvoyer.
   * @param entityId L'entityId du joueur à renvoyer.
   * @return Le joueur spécifié par le tick, ou null s'il n'existe pas.
   */
  public String getPlayerName(long tick, int entityId) {
    Snapshots<String> snapshots = map.get(entityId);
    if (snapshots == null) {
      return null;
    }
    return snapshots.getSnapshot(tick);
  }

  /**
   * Renvoit true si l'entité spécifiée par l'entityId est un joueur au tick spécifié.
   *
   * @param tick Le tick de l'entité.
   * @param entityId L'entityId de l'entité.
   * @return true si l'entité à ce tick est un joueur.
   */
  public boolean isPlayer(long tick, int entityId) {
    return getPlayerName(tick, entityId) != null;
  }

  /**
   * Enlève un joueur de la liste de joueurs, au tick spécifié.
   *
   * @param tick Le tick auquel enlever le joueur.
   * @param entityId L'entityId du joueur à enlever.
   */
  public void removePlayer(long tick, int entityId) {
    addPlayer(tick, entityId, null);
  }

  /**
   * Renvoit un IntStream des entityId des joueurs au tick spécifié.
   *
   * @param tick Le tick auquel récupérer la liste de joueurs.
   * @return Un stream d'entityId de joueurs.
   */
  public IntStream getPlayers(long tick) {
    return map.entrySet().stream().filter((e) -> e.getValue().getSnapshot(tick) != null).mapToInt((e) -> e.getKey());
  }

  /**
   * Indique au gestionnaire de joueurs que les joueurs appartenant à des ticks avant ou égaux au
   * tick spécifié ne seront plus jamais demandés et peuvent être supprimés.
   *
   * @param tick Le tick (inclus) jusqu'auquel les joueurs ne seront plus demandés.
   */
  public void disposeUntilTick(long tick) {
    // NOTE on pourrait ignorer ce call la plupart du temps pour éviter de parcourir toute la map à
    // chaque fois ; il faudrait sans doute modifier la capacité standard
    Set<Map.Entry<Integer, Snapshots<String>>> entrySet = map.entrySet();
    Iterator<Map.Entry<Integer, Snapshots<String>>> entryIterator = entrySet.iterator();
    while (entryIterator.hasNext()) {
      Map.Entry<Integer, Snapshots<String>> entry = entryIterator.next();
      Snapshots<String> snapshots = entry.getValue();
      snapshots.disposeUntilTick(tick);
      if (snapshots.isEmpty()) {
        entryIterator.remove();
      }
    }
  }

}
