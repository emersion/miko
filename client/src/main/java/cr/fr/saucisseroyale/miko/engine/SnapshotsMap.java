package cr.fr.saucisseroyale.miko.engine;

import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;

/**
 * Une map de snapshots, disposant de snapshots à chaque valeur de T. La lecture de la documentation
 * de {@link Snapshots} est recommandée.
 * <p>
 * Utile pour avoir un moyen de faire varier temporellement et indépendemment des valeurs
 * différentes selon un paramètre T.
 *
 * @param <T> Le type du paramètre pour obtenir une snapshot ("clef" de la map).
 * @param <U> Le type des valeurs du snapshot (variant temporellement).
 * @see Snapshots
 */
public class SnapshotsMap<T, U> {
  private static final int LINKED_LIST = -1;
  private final int snapshotsCapacity;
  private Map<T, Snapshots<U>> map = new HashMap<>();

  /**
   * Crée une map de snapshots, créant des snapshots avec {@link Snapshots#Snapshots()}.
   */
  public SnapshotsMap() {
    snapshotsCapacity = LINKED_LIST;
  }

  /**
   * Crée une map de snapshots, créant des snapshots avec {@link Snapshots#Snapshots(int)}.
   *
   * @param snapshotsCapacity La capacité initiale des snapshots à créer.
   */
  public SnapshotsMap(int snapshotsCapacity) {
    if (snapshotsCapacity <= 0) {
      throw new IllegalArgumentException("Snapshots capacity must be positive");
    }
    this.snapshotsCapacity = snapshotsCapacity;
  }

  /**
   * Ajoute une valeur à la map, à la clef spécifiée, au tick spécifié.
   *
   * @param tick  Le tick auquel ajouter la valeur.
   * @param key   La clef à laquelle ajouter la valeur.
   * @param value La valeur à ajouter.
   */
  public void setSnapshot(long tick, T key, U value) {
    Snapshots<U> snapshots = map.get(key);
    if (snapshots == null) {
      if (snapshotsCapacity == LINKED_LIST) {
        snapshots = new Snapshots<>();
      } else {
        snapshots = new Snapshots<>(snapshotsCapacity);
      }
      map.put(key, snapshots);
    }
    snapshots.setSnapshot(tick, value);
  }

  /**
   * Retourne la valeur à la clef spécifiée, au tick spécifié, ou null si cette valeur n'existe pas.
   *
   * @param tick Le tick de la valeur à renvoyer.
   * @param key  La clef de laquelle renvoyer la valeur.
   * @return La valeur spécifiée par la clef et le tick, ou null si elle n'existe pas.
   */
  public U getSnapshot(long tick, T key) {
    Snapshots<U> snapshots = map.get(key);
    if (snapshots == null) {
      return null;
    }
    return snapshots.getSnapshot(tick);
  }

  /**
   * Indique que toutes les valeurs appartenant à des ticks avant ou égaux au tick spécifié ne
   * seront plus jamais demandées et peuvent être supprimées.
   *
   * @param tick Le tick (inclus) jusqu'auquel les valeurs ne seront plus demandées.
   */
  public void disposeUntilTick(long tick) {
    Set<Entry<T, Snapshots<U>>> entrySet = map.entrySet();
    Iterator<Entry<T, Snapshots<U>>> entryIterator = entrySet.iterator();
    while (entryIterator.hasNext()) {
      Entry<T, Snapshots<U>> entry = entryIterator.next();
      Snapshots<U> snapshots = entry.getValue();
      snapshots.disposeUntilTick(tick);
      if (snapshots.isEmpty()) {
        entryIterator.remove();
      }
    }
  }
}
