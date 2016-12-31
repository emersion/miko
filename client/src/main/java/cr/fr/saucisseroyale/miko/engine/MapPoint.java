package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.TerrainPoint;
import cr.fr.saucisseroyale.miko.util.Pair.FloatFloat;

/**
 * Un point immutable de la carte de jeu, en coordonnées à virgule flottante. Utilisé pour indiquer
 * un emplacement sans répéter (x,y).
 */
public final class MapPoint {
  private float x;
  private float y;

  public MapPoint(TerrainPoint terrainPoint) {
    this(terrainPoint.getX(), terrainPoint.getY());
  }

  /**
   * @param x La coordonnée x du point.
   * @param y La coordonnée y du point.
   */
  public MapPoint(float x, float y) {
    this.x = x;
    this.y = y;
  }

  /**
   * @return La coordonnée x du point.
   */
  public float getX() {
    return x;
  }

  /**
   * @return La coordonnée y du point.
   */
  public float getY() {
    return y;
  }

  /**
   * Retourne ce point après translation de (deltaX;deltaY) (déplacé de deltaX en abscisse et de
   * deltaY en ordonnée).
   *
   * @param deltaX La translation en X à effectuer.
   * @param deltaY La translation en Y à effectuer.
   * @return Le point après translation.
   */
  public MapPoint getTranslated(float deltaX, float deltaY) {
    return new MapPoint(x + deltaX, y + deltaY);
  }

  public TerrainPoint toTerrainPoint() {
    return new TerrainPoint((int) x, (int) y);
  }

  public FloatFloat subtract(MapPoint other) {
    float deltaX = x - other.x;
    float deltaY = y - other.y;
    return new FloatFloat(deltaX, deltaY);
  }

  @Override
  public int hashCode() {
    final int prime = 31;
    int result = 1;
    result = prime * result + Float.floatToIntBits(x);
    result = prime * result + Float.floatToIntBits(y);
    return result;
  }

  @Override
  public boolean equals(Object obj) {
    if (this == obj) {
      return true;
    }
    if (obj == null) {
      return false;
    }
    if (!(obj instanceof MapPoint)) {
      return false;
    }
    MapPoint other = (MapPoint) obj;
    if (Float.floatToIntBits(x) != Float.floatToIntBits(other.x)) {
      return false;
    }
    return Float.floatToIntBits(y) == Float.floatToIntBits(other.y);
  }
}
