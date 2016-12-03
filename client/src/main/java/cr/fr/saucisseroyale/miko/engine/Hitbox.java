package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.Pair.DoubleFloat;
import cr.fr.saucisseroyale.miko.util.SingleElementIterable;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Une hitbox d'entité.
 */
public abstract class Hitbox {
  /**
   * Renvoit true si les deux hitbox, placées respectivement aux coordoonées spécifiées,
   * s'intersectionnent ou sont contenues l'une dans l'autre. Cette fonction est symétrique.
   *
   * @param h1 La première hitbox.
   * @param m1 La position du centre de la première hitbox.
   * @param h2 La deuxième hitbox.
   * @param m2 La position du centre de la deuxième hitbox.
   * @return true si les hitbox placées aux coordonnées spécifiées s'intersectionnent.
   */
  public static boolean collide(Hitbox h1, MapPoint m1, Hitbox h2, MapPoint m2) {
    if (h1 instanceof Null || h2 instanceof Null) {
      return false;
    }
    if (h1 instanceof Rectangle) {
      if (h2 instanceof Rectangle) {
        return rectangleCollide((Rectangle) h1, m1, (Rectangle) h2, m2);
      } else if (h2 instanceof Circle) {
        return rectCircleCollide((Rectangle) h1, m1, (Circle) h2, m2);
      } else {
        throw new IllegalArgumentException("Unknown hitbox type for h2 : " + h2);
      }
    } else if (h1 instanceof Circle) {
      if (h2 instanceof Rectangle) {
        return rectCircleCollide((Rectangle) h2, m2, (Circle) h1, m1);
      } else if (h2 instanceof Circle) {
        return circleCollide((Circle) h1, m1, (Circle) h2, m2);
      } else {
        throw new IllegalArgumentException("Unknown hitbox type for h2 : " + h2);
      }
    } else {
      throw new IllegalArgumentException("Unknown hitbox type for h1 : " + h1);
    }
  }

  public static Hitbox newCircleHitbox(float radius) {
    return new Circle(radius);
  }

  public static Hitbox newRectangleHitbox(float width, float height) {
    return new Rectangle(width, height);
  }

  public static Hitbox newNullHitbox() {
    return new Null();
  }

  private static boolean circleCollide(Circle h1, MapPoint m1, Circle h2, MapPoint m2) {
    float deltaX = m1.getX() - m2.getX();
    float deltaY = m1.getY() - m2.getY();
    float radiusSum = h1.radius + h2.radius;
    return deltaX * deltaX + deltaY * deltaY <= radiusSum * radiusSum;
  }

  private static boolean rectCircleCollide(Rectangle h1, MapPoint m1, Circle h2, MapPoint m2) {
    // taken from http://stackoverflow.com/a/402010
    float deltaX = Math.abs(m2.getX() - m1.getX());
    float deltaY = Math.abs(m2.getY() - m1.getY());
    if (deltaX > h1.widthHalf + h2.radius) {
      return false;
    }
    if (deltaY > h1.heightHalf + h2.radius) {
      return false;
    }
    if (deltaX <= h1.widthHalf) {
      return true;
    }
    if (deltaY <= h1.heightHalf) {
      return true;
    }
    float cornerX = deltaX - h1.widthHalf;
    float cornerY = deltaY - h1.heightHalf;
    return cornerX * cornerX + cornerY * cornerY <= h2.radius * h2.radius;
  }

  private static boolean rectangleCollide(Rectangle h1, MapPoint m1, Rectangle h2, MapPoint m2) {
    MapPoint upLeft1 = m1.getTranslated(-h1.widthHalf, h1.heightHalf);
    MapPoint downRight1 = m1.getTranslated(h1.widthHalf, -h1.heightHalf);
    MapPoint upLeft2 = m2.getTranslated(-h2.widthHalf, h2.heightHalf);
    MapPoint downRight2 = m2.getTranslated(h2.widthHalf, -h2.heightHalf);
    if (upLeft1.getX() > downRight2.getX() || upLeft2.getX() > downRight1.getX()) {
      return false;
    }
    return !(upLeft1.getY() < downRight2.getY() || upLeft2.getY() < downRight1.getY());
  }

  /**
   * Retourne les points symbolisant le "contour utile" lors d'une translation de (deltaX;deltaY).
   * Lors de la translation, on vérifiera les collisions de ce points avec le terrain. C'est une
   * approximation du contour réel, qui dépend de la direction (et éventuellement de la norme) de la
   * translation.
   *
   * @param deltaX La différence d'abscisses de la translation.
   * @param deltaY La différence d'ordonnées de la translation.
   * @return Les points du contour utile, en paires de décalage en x et y par rapport au centre de
   * la hitbox/l'entité.
   */
  public abstract Iterable<DoubleFloat> getKeyPoints(float deltaX, float deltaY);

  private static class Null extends Hitbox {
    @Override
    public Iterable<DoubleFloat> getKeyPoints(float deltaX, float deltaY) {
      return Collections.emptyList();
    }
  }

  private static class Circle extends Hitbox {
    private final float radius;
    private final List<DoubleFloat> keyPoints;

    public Circle(float radius) {
      this.radius = radius;
      // for now approximate circle to ♦
      // TODO might be optimized to only 4 points
      List<DoubleFloat> pairList = new ArrayList<>(5);
      pairList.add(new DoubleFloat(0f, 0f));
      pairList.add(new DoubleFloat(-radius, 0f));
      pairList.add(new DoubleFloat(radius, 0f));
      pairList.add(new DoubleFloat(0f, -radius));
      pairList.add(new DoubleFloat(0f, radius));
      keyPoints = Collections.unmodifiableList(pairList);
    }

    @Override
    public Iterable<DoubleFloat> getKeyPoints(float deltaX, float deltaY) {
      return keyPoints;
    }
  }

  private static class Rectangle extends Hitbox {
    private final float widthHalf;
    private final float heightHalf;
    private final List<DoubleFloat> keyPoints;

    public Rectangle(float width, float height) {
      widthHalf = width / 2;
      heightHalf = height / 2;
      List<DoubleFloat> pairs = new ArrayList<>(5);
      pairs.add(new DoubleFloat(0f, 0f));
      pairs.add(new DoubleFloat(widthHalf, heightHalf));
      pairs.add(new DoubleFloat(widthHalf, -heightHalf));
      pairs.add(new DoubleFloat(-widthHalf, heightHalf));
      pairs.add(new DoubleFloat(-widthHalf, -heightHalf));
      keyPoints = Collections.unmodifiableList(pairs);
    }

    @Override
    public Iterable<DoubleFloat> getKeyPoints(float deltaX, float deltaY) {
      if (deltaX == 0) {
        if (deltaY > 0) {
          return new SingleElementIterable<>(new DoubleFloat(0f, heightHalf));
        } else if (deltaY < 0) {
          return new SingleElementIterable<>(new DoubleFloat(0f, -heightHalf));
        } else {
          return Collections.emptyList();
        }
      }
      if (deltaY == 0) {
        if (deltaX > 0) {
          return new SingleElementIterable<>(new DoubleFloat(widthHalf, 0f));
        } else {
          return new SingleElementIterable<>(new DoubleFloat(-widthHalf, 0f));
        }
      }
      int excludeIndex;
      if (deltaX > 0) {
        excludeIndex = deltaY > 0 ? 1 : 2;
      } else {
        excludeIndex = deltaY > 0 ? 3 : 4;
      }
      List<DoubleFloat> points = new ArrayList<>(4);
      for (int i = 0, n = keyPoints.size(); i < n; i++) {
        if (i != excludeIndex) {
          points.add(keyPoints.get(i));
        }
      }
      return points;
    }
  }
}
