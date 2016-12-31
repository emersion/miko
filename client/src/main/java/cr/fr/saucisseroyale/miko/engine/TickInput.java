package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.MikoMath;
import cr.fr.saucisseroyale.miko.util.Or;
import cr.fr.saucisseroyale.miko.util.Pair;
import fr.delthas.uitest.Ui;

import java.awt.geom.Point2D;
import java.util.Collections;
import java.util.List;
import java.util.prefs.Preferences;

/**
 * Un input pour un tick d'engine de Miko.
 */
class TickInput {
  private static final Preferences prefsNode = Preferences.userRoot().node("miko.input");
  private static final int moveLeftKeycode = prefsNode.getInt("left", Ui.KEY_A);
  private static final int moveRightKeycode = prefsNode.getInt("right", Ui.KEY_D);
  private static final int moveUpKeycode = prefsNode.getInt("up", Ui.KEY_W);
  private static final int moveDownKeycode = prefsNode.getInt("down", Ui.KEY_S);
  private static final int ballSendKeycode = prefsNode.getInt("ball", Ui.MOUSE_LEFT);
  private boolean moveLeft;
  private boolean moveRight;
  private boolean moveUp;
  private boolean moveDown;
  private int lastLeftRightPressed = 0;
  private int lastUpDownPressed = 0;
  private boolean ballSendDown;
  private float ballSendAngle = Float.NaN;

  public TickInput(TickInput previous, List<Or<Pair.DoubleDouble, Pair.IntBoolean>> eventList) {

    if (previous != null) {
      moveLeft = previous.moveLeft;
      moveRight = previous.moveRight;
      moveUp = previous.moveUp;
      moveDown = previous.moveDown;
      lastLeftRightPressed = previous.lastLeftRightPressed;
      lastUpDownPressed = previous.lastUpDownPressed;
      ballSendDown = previous.ballSendDown;
    }

    float mouseAngle = Float.NaN;

    for (Or<Pair.DoubleDouble, Pair.IntBoolean> event : eventList) {
      if (event.isFirst()) {
        mouseAngle = getAngleFromMousePosition(event.getAsFirst());
        continue;
      }
      int key = event.getAsSecond().getFirst();
      if (event.getAsSecond().getSecond()) {
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
          break;
        }
        if (key == moveDownKeycode) {
          lastUpDownPressed = 1;
          moveDown = true;
          continue;
        }
        if (key == ballSendKeycode) {
          ballSendDown = true;
          ballSendAngle = mouseAngle;
          continue;
        }
      } else {
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
        if (key == ballSendKeycode) {
          ballSendDown = false;
          continue;
        }
      }
    }

    if (ballSendDown && Float.isNaN(ballSendAngle)) {
      ballSendAngle = mouseAngle;
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

  private static float getAngleFromMousePosition(Pair.DoubleDouble mousePosition) {
    if (mousePosition == null) {
      return Float.NaN;
    }
    // angle can be calculated from screen middle because the render translation (offset based
    // on mouse position) has the same direction
    float angle = MikoMath.atan2(new Point2D.Double(mousePosition.getFirst() - Ui.getWidth() / 2.0, mousePosition.getSecond() - Ui.getHeight() / 2.0));
    return angle;
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
