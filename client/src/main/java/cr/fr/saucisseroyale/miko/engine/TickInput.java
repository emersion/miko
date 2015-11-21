package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.MikoMath;
import cr.fr.saucisseroyale.miko.util.Triplet;

import java.awt.Point;
import java.awt.event.KeyEvent;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.prefs.Preferences;

/**
 * Un input pour un tick d'engine de Miko.
 *
 */
class TickInput {

  private static final Preferences prefsNode = Preferences.userRoot().node("miko.input");

  private static final int ballSendKeycode = prefsNode.getInt("ballSend", KeyEvent.VK_SPACE);
  private static final int moveLeftKeycode = prefsNode.getInt("left", KeyEvent.VK_Q);
  private static final int moveRightKeycode = prefsNode.getInt("right", KeyEvent.VK_D);
  private static final int moveUpKeycode = prefsNode.getInt("up", KeyEvent.VK_Z);
  private static final int moveDownKeycode = prefsNode.getInt("down", KeyEvent.VK_S);

  private boolean moveLeft;
  private boolean moveRight;
  private boolean moveUp;
  private boolean moveDown;

  // list of ball send requests represented by send direction
  private List<Float> ballSendRequests = new LinkedList<>();

  public TickInput(TickInput previous, List<Triplet<Boolean, Integer, Point>> eventList) {

    int lastLeftRightPressed = 0;
    int lastUpDownPressed = 0;

    boolean newLeft = false, newRight = false, newUp = false, newDown = false;

    for (Triplet<Boolean, Integer, Point> event : eventList) {
      int key = event.getSecond();
      if (event.getFirst()) {
        if (key == ballSendKeycode) {
          // angle can be calculated from screen middle because the render translation (offset based
          // on mouse position) has the same direction
          Point mousePosition = event.getThird();
          if (mousePosition == null) {
            continue;
          }
          float angle = MikoMath.atan2(mousePosition);
          ballSendRequests.add(angle);
        }
        if (key == moveLeftKeycode) {
          lastLeftRightPressed = -1;
          newLeft = true;
          continue;
        }
        if (key == moveRightKeycode) {
          lastLeftRightPressed = 1;
          newRight = true;
          continue;
        }
        if (key == moveUpKeycode) {
          lastUpDownPressed = -1;
          newUp = true;
          continue;
        }
        if (key == moveDownKeycode) {
          lastUpDownPressed = 1;
          newDown = true;
          continue;
        }
      } else {
        if (key == moveLeftKeycode) {
          newLeft = false;
          continue;
        }
        if (key == moveRightKeycode) {
          newRight = false;
          continue;
        }
        if (key == moveUpKeycode) {
          newUp = false;
          continue;
        }
        if (key == moveDownKeycode) {
          newDown = false;
          continue;
        }
      }
    }

    if (newLeft && newRight) {
      moveLeft = lastLeftRightPressed < 0;
      moveRight = lastLeftRightPressed > 0;
    } else {
      moveLeft = newLeft;
      moveRight = newRight;
    }

    if (newUp && newDown) {
      moveUp = lastUpDownPressed < 0;
      moveDown = lastUpDownPressed > 0;
    } else {
      moveUp = newUp;
      moveDown = newDown;
    }

  }

  /**
   * @return true si la touche "se déplacer vers la gauche" est appuyée.
   */
  public boolean isMoveLeft() {
    return moveLeft;
  }

  /**
   * @return true si la touche "se déplacer vers la droite" est appuyée.
   */
  public boolean isMoveRight() {
    return moveRight;
  }

  /**
   * @return true si la touche "se déplacer vers le haut" est appuyée.
   */
  public boolean isMoveUp() {
    return moveUp;
  }

  /**
   * @return true si la touche "se déplacer vers le bas" est appuyée.
   */
  public boolean isMoveDown() {
    return moveDown;
  }

  /**
   * @return La liste des requêtes de lancement de balles.
   */
  public Iterable<Float> getBallSendRequests() {
    return ballSendRequests;
  }

  /**
   * Retourne l'input par défaut suivant un input donné.
   *
   * @param previous L'input sur lequel se baser pour construire l'input.
   * @return Un input par défaut basé sur l'input passé en paramètre.
   */
  public static TickInput getNextFrom(TickInput previous) {
    return new TickInput(previous, Collections.emptyList());
  }

}
