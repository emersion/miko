package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.Or;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.*;
import java.util.Map.Entry;

/**
 * Un gestionnaire des inputs à chaque tick de jeu.
 */
class TickInputManager {
  private Map<Long, TickInput> map = new HashMap<>();
  private long firstTick = Long.MAX_VALUE;

  public void addInput(long tick, List<Or<Pair.DoubleDouble, Pair.IntBoolean>> eventList) {
    TickInput previousInput = getInput(tick - 1);
    TickInput newInput = new TickInput(previousInput, eventList);
    map.put(tick, newInput);
    if (tick < firstTick) {
      firstTick = tick;
    }
  }

  public TickInput getInput(long tick) {
    if (tick < firstTick) {
      return null;
    }
    TickInput tickInput = map.get(tick);
    if (tickInput == null) {
      TickInput previousInput = getInput(tick - 1);
      TickInput newInput = TickInput.getNextFrom(previousInput);
      map.put(tick, newInput);
      return newInput;
    } else {
      return tickInput;
    }
  }

  public void disposeUntilTick(long tick) {
    Set<Entry<Long, TickInput>> entrySet = map.entrySet();
    Iterator<Entry<Long, TickInput>> iterator = entrySet.iterator();
    while (iterator.hasNext()) {
      Entry<Long, TickInput> entry = iterator.next();
      if (entry.getKey() > tick) {
        firstTick = entry.getKey();
        break;
      }
      iterator.remove();
    }
  }
}
