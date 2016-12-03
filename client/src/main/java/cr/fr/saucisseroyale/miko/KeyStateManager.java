package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.util.Triplet;

import java.awt.*;
import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;

/**
 * Un gestionnaire de l'état des touches du clavier, disposant des touches pressées (avec la
 * position de la souris) depuis le dernier flush.
 */
class KeyStateManager implements KeyListener {
  private final Supplier<Point> mousePositionSupplier;
  // additional synchronization to ensure no event gets lost
  private final Object eventListLock = new Object();
  private List<Triplet<Boolean, Integer, Point>> eventList = new ArrayList<>();

  public KeyStateManager(Supplier<Point> mousePositionSupplier) {
    this.mousePositionSupplier = mousePositionSupplier;
  }

  public List<Triplet<Boolean, Integer, Point>> getEventsAndFlush() {
    List<Triplet<Boolean, Integer, Point>> lastFrameList = eventList;
    synchronized (eventListLock) {
      eventList = new ArrayList<>();
    }
    return lastFrameList;
  }

  @Override
  public void keyTyped(KeyEvent e) {
    // ignore type events
  }

  @Override
  public void keyPressed(KeyEvent e) {
    synchronized (eventListLock) {
      eventList.add(new Triplet<>(Boolean.TRUE, e.getKeyCode(), mousePositionSupplier.get()));
    }
  }

  @Override
  public void keyReleased(KeyEvent e) {
    synchronized (eventListLock) {
      eventList.add(new Triplet<>(Boolean.FALSE, e.getKeyCode(), mousePositionSupplier.get()));
    }
  }
}
