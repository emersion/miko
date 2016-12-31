package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.util.Or;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.ArrayList;
import java.util.List;

/**
 * Un gestionnaire de l'état des touches du clavier, disposant des touches pressées (avec la
 * position de la souris) depuis le dernier flush.
 */
class InputStateManager {
  private List<Or<Pair.DoubleDouble, Pair.IntBoolean>> eventList = new ArrayList<>();
  private double mouseX;
  private double mouseY;

  public List<Or<Pair.DoubleDouble, Pair.IntBoolean>> getEventsAndFlush() {
    List<Or<Pair.DoubleDouble, Pair.IntBoolean>> lastFrameList = eventList;
    eventList = new ArrayList<>();
    eventList.add(new Or<>(new Pair.DoubleDouble(mouseX, mouseY), true));
    return lastFrameList;
  }

  public void pushKeyButton(int key, boolean down) {
    eventList.add(new Or<>(new Pair.IntBoolean(key, down), false));
  }

  public void pushMouseButton(int button, boolean down) {
    eventList.add(new Or<>(new Pair.IntBoolean(button, down), false));
  }

  public void pushMouseMove(double x, double y) {
    mouseX = x;
    mouseY = y;
    eventList.add(new Or<>(new Pair.DoubleDouble(mouseX, mouseY), true));
  }
}
