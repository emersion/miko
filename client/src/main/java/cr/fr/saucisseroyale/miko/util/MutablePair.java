package cr.fr.saucisseroyale.miko.util;

/**
 * A mutable generic pair of data (2-tuple).
 *
 * @param <T> The type of the first element.
 * @param <U> The type of the second element.
 *
 */
public class MutablePair<T, U> {

  private T first;
  private U second;

  /**
   * @param first The first element of the pair.
   * @param second The second element of the pair.
   */
  public MutablePair(T first, U second) {
    this.first = first;
    this.second = second;
  }

  /**
   * @return The first element of the pair.
   */
  public T getFirst() {
    return first;
  }

  /**
   * @return The second element of the pair.
   */
  public U getSecond() {
    return second;
  }

  /**
   * @param first The first element of the pair to set.
   */
  public void setFirst(T first) {
    this.first = first;
  }

  /**
   * @param second The second element of the pair to set.
   */
  public void setSecond(U second) {
    this.second = second;
  }

}
