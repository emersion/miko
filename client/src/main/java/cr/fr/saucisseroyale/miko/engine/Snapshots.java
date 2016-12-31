package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.MutablePair;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.ListIterator;

/**
 * Un gestionnaire de snapshots, stockant des couples (tick ; donnée) et renvoyant pour un tick
 * demandé, la donnée la plus récente dont le tick est avant celui demandé.
 * <p>
 * L'ajout de snapshots est effectué grâce à {@link #setSnapshot(long, Object)}, la lecture grâce à
 * {@link #getSnapshot(long)}.
 * <p>
 * Les utilisateurs de cette classe peuvent indiquer qu'ils n'utiliseront plus des valeurs
 * suffisament anciennes grâce à {@link #disposeUntilTick(long)}.
 *
 * @param <T> Le type de données à stocker dans chaque couple (tick ; donnée).
 */
class Snapshots<T> {
  private List<MutablePair.Long<T>> snapshots;
  private long first, last;

  /**
   * Crée un snapshots avec comme support une {@link LinkedList}.
   */
  public Snapshots() {
    snapshots = new LinkedList<>();
  }

  /**
   * Crée un snapshots avec comme support une {@link ArrayList} de capacité intiale spécifiée.
   *
   * @param capacity La capacité initiale de la {@link ArrayList} à créer.
   */
  public Snapshots(int capacity) {
    snapshots = new ArrayList<>(capacity);
  }

  /**
   * Retourne le snapshot au tick spécifié, ou null s'il n'existe pas de snapshot à ce tick.
   *
   * @param tick Le tick du chunk à renvoyer.
   * @return Le snapshot spécifié par le tick, ou null s'il n'existe pas.
   */
  public T getSnapshot(long tick) {
    if (tick - first < last - tick) {
      return getSnapshotFromBeginning(tick);
    } else {
      return getSnapshotFromEnd(tick);
    }
  }

  private T getSnapshotFromBeginning(long tick) {
    ListIterator<MutablePair.Long<T>> snapshotIterator = snapshots.listIterator();
    while (snapshotIterator.hasNext()) {
      MutablePair.Long<T> pair = snapshotIterator.next();
      if (pair.getFirst() == tick) {
        return pair.getSecond();
      }
      if (pair.getFirst() > tick) {
        snapshotIterator.previous();
        if (!snapshotIterator.hasPrevious()) {
          return null;
        }
        return snapshotIterator.previous().getSecond();
      }
    }
    if (snapshots.isEmpty()) {
      return null;
    }
    return snapshotIterator.previous().getSecond();
  }

  private T getSnapshotFromEnd(long tick) {
    ListIterator<MutablePair.Long<T>> snapshotIterator = snapshots.listIterator(snapshots.size());
    while (snapshotIterator.hasPrevious()) {
      MutablePair.Long<T> pair = snapshotIterator.previous();
      if (pair.getFirst() <= tick) {
        return pair.getSecond();
      }
    }
    return null;
  }

  /**
   * Ajoute un snapshot à la liste de snapshots, au tick spécifié.
   * <p>
   * Le snapshot sera aussi ajouté à tous les ticks suivants celui indiqué, jusqu'à ce qu'un autre
   * snapshot soit ajouté. Autrement dit, les ticks ne possédant pas de snapshots sont remplis avec
   * les snapshots les plus récents avant ce tick.
   *
   * @param tick     Le tick auquel ajouter le snapshot.
   * @param snapshot Le snapshot à ajouter.
   */
  public void setSnapshot(long tick, T snapshot) {
    if (tick - first < last - tick) {
      setSnapshotFromBeginning(tick, snapshot);
    } else {
      setSnapshotFromEnd(tick, snapshot);
    }
    if (tick < first) {
      first = tick;
    }
    if (tick > last) {
      last = tick;
    }
  }

  private void setSnapshotFromBeginning(long tick, T snapshot) {
    ListIterator<MutablePair.Long<T>> snapshotIterator = snapshots.listIterator();
    while (snapshotIterator.hasNext()) {
      MutablePair.Long<T> pair = snapshotIterator.next();
      if (pair.getFirst() == tick) {
        pair.setSecond(snapshot);
        return;
      }
      if (pair.getFirst() > tick) {
        snapshotIterator.previous();
        snapshotIterator.add(new MutablePair.Long<>(tick, snapshot));
        return;
      }
    }
    snapshots.add(new MutablePair.Long<>(tick, snapshot));
  }

  private void setSnapshotFromEnd(long tick, T snapshot) {
    ListIterator<MutablePair.Long<T>> snapshotIterator = snapshots.listIterator(snapshots.size());
    while (snapshotIterator.hasPrevious()) {
      MutablePair.Long<T> pair = snapshotIterator.previous();
      if (pair.getFirst() == tick) {
        pair.setSecond(snapshot);
        return;
      }
      if (pair.getFirst() < tick) {
        snapshotIterator.next();
        snapshotIterator.add(new MutablePair.Long<>(tick, snapshot));
        return;
      }
    }
    snapshots.add(new MutablePair.Long<>(tick, snapshot));
  }

  /**
   * Indique que les snapshots appartenant à des ticks avant ou égaux au tick spécifié ne seront
   * plus jamais demandés et peuvent être supprimés.
   * <p>
   * Cette méthode n'est qu'une indication ; la classe peut décider de disposer les ressources
   * associées aux ticks, ou non.
   *
   * @param tick Le tick (inclus) jusqu'auquel les snapshots ne seront plus demandés.
   */
  public void disposeUntilTick(long tick) {
    if (snapshots.isEmpty()) {
      return;
    }
    ListIterator<MutablePair.Long<T>> snapshotIterator = snapshots.listIterator();
    while (snapshotIterator.hasNext()) {
      MutablePair.Long<T> pair = snapshotIterator.next();
      if (pair.getFirst() > tick) {
        snapshotIterator.previous();
        break;
      }
    }
    if (!snapshotIterator.hasPrevious()) {
      return;
    }
    first = snapshotIterator.previous().getFirst();
    while (snapshotIterator.hasPrevious()) {
      snapshotIterator.previous();
      snapshotIterator.remove();
    }
  }

  /**
   * Renvoit false si la liste de snapshots n'est pas vide, et peut renvoyer true si elle est vide.
   * Méthode de convenance retournant exactement {@code size()==0}.
   *
   * @return Le booléen décrit ci-dessus.
   */
  public boolean isEmpty() {
    return size() == 0;
  }

  /**
   * Renvoit une valeur supérieur ou égale au nombre de snapshots enregistrés dans la liste. Appeler
   * {@link #disposeUntilTick(long)} puis cette méthode <b>ne renverra pas forcément</b> le nombre
   * effectif de snapshots après le tick.
   *
   * @return Une majoration du nombre de snapshots stockés.
   */
  public int size() {
    return snapshots.size();
  }
}
