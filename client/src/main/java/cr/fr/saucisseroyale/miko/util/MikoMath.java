package cr.fr.saucisseroyale.miko.util;

import java.awt.Point;

/**
 * Classe d'utilité avec des implémentations personnalisées de fonctions mathématiques.
 *
 */
public class MikoMath {

  // classe statique
  private MikoMath() {
    throw new IllegalStateException("This class can't be instantiated");
  }

  public static int quotient(int number, int divisor) {
    if (number < 0) {
      return -quotient(-number, divisor) - 1;
    }
    return number / divisor;
  }

  public static int modulo(int number, int divisor) {
    if (number < 0) {
      return modulo(-number, divisor);
    }
    return number % divisor;
  }

  public static float cos(float number) {
    return (float) Math.cos(number);
  }

  public static float sin(float number) {
    return (float) Math.sin(number);
  }

  public static float atan2(Point point) {
    return atan2(point.y, point.x);
  }

  public static float atan2(float y, float x) {
    return (float) Math.atan2(y, x);
  }

}
