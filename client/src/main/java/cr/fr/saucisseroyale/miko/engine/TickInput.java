package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.MikoMath;
import cr.fr.saucisseroyale.miko.util.Triplet;

import java.awt.*;
import java.awt.event.KeyEvent;
import java.util.Collections;
import java.util.List;
import java.util.prefs.Preferences;

/**
 * Un input pour un tick d'engine de Miko.
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
  private float ballSendAngle = Float.NaN;

  public TickInput(TickInput previous, List<Triplet<Boolean, Integer, Point>> eventList, Point mousePosition) {

    int lastLeftRightPressed = 0;
    int lastUpDownPressed = 0;

    if (previous != null) {
      moveLeft = previous.moveLeft;
      moveRight = previous.moveRight;
      moveUp = previous.moveUp;
      moveDown = previous.moveDown;
      if (!Float.isNaN(previous.ballSendAngle)) {
        ballSendAngle = getAngleFromMousePosition(mousePosition);
      }
    }

    for (Triplet<Boolean, Integer, Point> event : eventList) {
      int key = event.getSecond();
      if (event.getFirst()) {
        if (key == ballSendKeycode) {
          // angle can be calculated from screen middle because the render translation (offset based
          // on mouse position) has the same direction
          ballSendAngle = getAngleFromMousePosition(event.getThird());
        }
        if (key == moveLeftKeycode) {
          lastLeftRightPressed = -1;
          moveLeft = true;
          continue;
        }
        if (key == moveRightKeycode) {
          lastLeftRightPressed = 1;
          moveRight = true;
          continue;
        }
        if (key == moveUpKeycode) {
          lastUpDownPressed = -1;
          moveUp = true;
          continue;
        }
        if (key == moveDownKeycode) {
          lastUpDownPressed = 1;
          moveDown = true;
          continue;
        }
      } else {
        if (key == ballSendKeycode) {
          ballSendAngle = Float.NaN;
        }
        if (key == moveLeftKeycode) {
          moveLeft = false;
          continue;
        }
        if (key == moveRightKeycode) {
          moveRight = false;
          continue;
        }
        if (key == moveUpKeycode) {
          moveUp = false;
          continue;
        }
        if (key == moveDownKeycode) {
          moveDown = false;
          continue;
        }
      }
    }

    if (moveLeft && moveRight) {
      moveLeft = lastLeftRightPressed < 0;
      moveRight = lastLeftRightPressed > 0;
    }

    if (moveUp && moveDown) {
      moveUp = lastUpDownPressed < 0;
      moveDown = lastUpDownPressed > 0;
    }
  }

  private static float getAngleFromMousePosition(Point mousePosition) {
    if (mousePosition == null) {
      return Float.NaN;
    }
    // angle can be calculated from screen middle because the render translation (offset based
    // on mouse position) has the same direction
    float angle = MikoMath.atan2(mousePosition);
    return angle;
  }

  /**
   * Retourne l'input par défaut suivant un input donné.
   *
   * @param previous L'input sur lequel se baser pour construire l'input.
   * @return Un input par défaut basé sur l'input passé en paramètre.
   */
  public static TickInput getNextFrom(TickInput previous) {
    return new TickInput(previous, Collections.emptyList(), null);
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
   * Renvoit une requête de lancement de balle sous la forme de l'angle de lancement désiré. Si le
   * joueur n'a pas lancé de balle, renvoit Float.NaN.
   *
   * @return L'angle de lancement de balle souhaité, ou Float.NaN.
   */
  public float getBallSendRequest() {
    return ballSendAngle;
  }
}
