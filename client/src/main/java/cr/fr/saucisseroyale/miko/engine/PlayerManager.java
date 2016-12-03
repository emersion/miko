package cr.fr.saucisseroyale.miko.engine;

/**
 * Un gestionnaire de joueurs, stockant tous les joueurs à tous les ticks, avec le principe de
 * {@link Snapshots}.
 *
 * @see Snapshots
 */
class PlayerManager {
  private SnapshotsMap<Integer, String> map = new SnapshotsMap<>();

  /**
   * Ajoute un joueur à la liste de joueurs, au tick spécifié.
   *
   * @param tick     Le tick auquel ajouter le joueur.
   * @param entityId L'entityId du joueur à ajouter.
   * @param name     Le nom du joueur à ajouter.
   */
  public void addPlayer(long tick, int entityId, String name) {
    map.setSnapshot(tick, entityId, name);
  }

  /**
   * Retourne le joueur au tick spécifié, ou null s'il n'existe pas de joueur à ce tick.
   *
   * @param tick     Le tick du joueur à renvoyer.
   * @param entityId L'entityId du joueur à renvoyer.
   * @return Le joueur spécifié par le tick, ou null s'il n'existe pas.
   */
  public String getPlayerName(long tick, int entityId) {
    return map.getSnapshot(tick, entityId);
  }

  /**
   * Renvoit true si l'entité spécifiée par l'entityId est un joueur au tick spécifié.
   *
   * @param tick     Le tick de l'entité.
   * @param entityId L'entityId de l'entité.
   * @return true si l'entité à ce tick est un joueur.
   */
  public boolean isPlayer(long tick, int entityId) {
    return getPlayerName(tick, entityId) != null;
  }

  /**
   * Enlève un joueur de la liste de joueurs, au tick spécifié.
   *
   * @param tick     Le tick auquel enlever le joueur.
   * @param entityId L'entityId du joueur à enlever.
   */
  public void removePlayer(long tick, int entityId) {
    addPlayer(tick, entityId, null);
  }

  /**
   * Indique au gestionnaire de joueurs que les joueurs appartenant à des ticks avant ou égaux au
   * tick spécifié ne seront plus jamais demandés et peuvent être supprimés.
   *
   * @param tick Le tick (inclus) jusqu'auquel les joueurs ne seront plus demandés.
   */
  public void disposeUntilTick(long tick) {
    map.disposeUntilTick(tick);
  }
}
