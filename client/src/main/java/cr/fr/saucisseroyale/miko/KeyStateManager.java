package cr.fr.saucisseroyale.miko;

import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;
import java.util.Comparator;
import java.util.HashMap;
import java.util.Map;
import java.util.stream.IntStream;

/**
 * Un gestionnaire de l'état des touches du clavier, disposant des touches actuellement pressées, ou
 * pressées depuis le dernier flush.
 *
 */
class KeyStateManager implements KeyListener {

  private Map<Integer, Long> keyMap = new HashMap<>();
  private long lastFlushTime = 0;

  public IntStream getPressedKeys() {
    return keyMap.entrySet().stream().filter((e) -> (e.getValue() > 0)).mapToInt((e) -> e.getKey());
  }

  public IntStream flush() {
    long time = System.nanoTime();
    IntStream stream =
        keyMap.entrySet().stream().filter((e) -> (Math.abs(e.getValue()) > lastFlushTime))
        .sorted(Comparator.comparingLong((e) -> Math.abs(e.getValue()))).mapToInt((e) -> e.getKey());
    lastFlushTime = time;
    return stream;
  }

  @Override
  public void keyTyped(KeyEvent e) {
    // ignore type events
  }

  @Override
  public void keyPressed(KeyEvent e) {
    long time = System.nanoTime();
    keyMap.put(e.getKeyCode(), time);
  }

  @Override
  public void keyReleased(KeyEvent e) {
    Long time = keyMap.get(e.getKeyCode());
    if (time == null || time <= 0) {
      return;
    }
    keyMap.put(e.getKeyCode(), -time);
  }

}
