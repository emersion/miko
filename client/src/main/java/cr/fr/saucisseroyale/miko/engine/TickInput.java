package cr.fr.saucisseroyale.miko.engine;

import java.awt.Point;
import java.awt.event.KeyEvent;
import java.util.prefs.Preferences;
import java.util.stream.IntStream;

/**
 * Un input pour un tick d'engine de Miko.
 *
 */
class TickInput {

  private static final Preferences prefsNode = Preferences.userRoot().node("miko.input");

  private static final int moveLeftKeycode = prefsNode.getInt("left", KeyEvent.VK_Q);
  private static final int moveRightKeycode = prefsNode.getInt("right", KeyEvent.VK_D);
  private static final int moveUpKeycode = prefsNode.getInt("up", KeyEvent.VK_Z);
  private static final int moveDownKeycode = prefsNode.getInt("down", KeyEvent.VK_S);

  private boolean moveLeft;
  private boolean moveRight;
  private boolean moveUp;
  private boolean moveDown;

  private final boolean mouseOver;
  private final int mouseX;
  private final int mouseY;

  public TickInput(TickInput previous, IntStream pressedKeys, IntStream newlyPressedKeys, Point mousePosition) {
    pressedKeys.forEach((keyCode) -> {
      if (keyCode == moveLeftKeycode) {
        moveLeft = true;
        return;
      }
      if (keyCode == moveRightKeycode) {
        moveRight = true;
        return;
      }
      if (keyCode == moveUpKeycode) {
        moveUp = true;
        return;
      }
      if (keyCode == moveDownKeycode) {
        moveDown = true;
        return;
      }
    });

    newlyPressedKeys.forEach((keyCode) -> {

    });

    if (mousePosition != null) {
      mouseOver = true;
      mouseX = mousePosition.x;
      mouseY = mousePosition.y;
    } else {
      mouseOver = false;
      if (previous != null) {
        mouseX = previous.mouseX;
        mouseY = previous.mouseY;
      } else {
        mouseX = 0;
        mouseY = 0;
      }
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
   * @return true si la souris est au-dessus du composant de jeu.
   */
  public boolean isMouseOver() {
    return mouseOver;
  }

  /**
   * @return La position en X de la souris.
   */
  public int getMouseX() {
    return mouseX;
  }

  /**
   * @return La position en Y de la souris.
   */
  public int getMouseY() {
    return mouseY;
  }

  /**
   * Retourne l'input par défaut suivant un input donné.
   *
   * @param previous L'input sur lequel se baser pour construire l'input.
   * @return Un input par défaut basé sur l'input passé en paramètre.
   */
  public static TickInput getNextFrom(TickInput previous) {
    return new TickInput(previous, IntStream.empty(), IntStream.empty(), null);
  }

}
