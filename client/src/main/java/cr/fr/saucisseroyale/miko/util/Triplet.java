package cr.fr.saucisseroyale.miko.util;

/**
 * An immutable generic triplet of data (3-tuple).
 * 
 * @param <T> The type of the first element.
 * @param <U> The type of the second element.
 * @param <V> The type of the third element.
 *
 */
public class Triplet<T, U, V> {

  private final T first;
  private final U second;
  private final V third;

  /**
   * @param first The first element of the triplet.
   * @param second The second element of the triplet.
   * @param third The third element of the triplet.
   */
  public Triplet(T first, U second, V third) {
    this.first = first;
    this.second = second;
    this.third = third;
  }

  /**
   * @return The first element of the triplet.
   */
  public T getFirst() {
    return first;
  }

  /**
   * @return The second element of the triplet.
   */
  public U getSecond() {
    return second;
  }

  /**
   * @return The third element of the triplet.
   */
  public V getThird() {
    return third;
  }



}
