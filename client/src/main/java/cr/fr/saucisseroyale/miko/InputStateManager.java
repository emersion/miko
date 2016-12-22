package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.util.Triplet;

import java.awt.*;
import java.util.ArrayList;
import java.util.List;

/**
 * Un gestionnaire de l'état des touches du clavier, disposant des touches pressées (avec la
 * position de la souris) depuis le dernier flush.
 */
class InputStateManager {
  private List<Triplet<Integer, Integer, Point.Double>> eventList = new ArrayList<>();
  private double mouseX;
  private double mouseY;

  public List<Triplet<Integer, Integer, Point.Double>> getEventsAndFlush() {
    List<Triplet<Integer, Integer, Point.Double>> lastFrameList = eventList;
    eventList = new ArrayList<>();
    eventList.add(new Triplet<>(0, 0, new Point.Double(mouseX, mouseY)));
    return lastFrameList;
  }

  public void pushKeyButton(double x, double y, int key, boolean down) {
    eventList.add(new Triplet<>(down ? 1 : 2, key, new Point.Double(x, y)));
  }

  public void pushMouseButton(double x, double y, int button, boolean down) {
    eventList.add(new Triplet<>(down ? 3 : 4, button, new Point.Double(x, y)));
  }

  public void pushMouseMove(double x, double y) {
    mouseX = x;
    mouseY = y;
    eventList.add(new Triplet<>(0, 0, new Point.Double(x, y)));
  }
}
