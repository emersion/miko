package cr.fr.saucisseroyale.miko.util;

/**
 * Un triplet générique immutable de trois éléments (3-tuple).
 *
 * @param <T> Le type du premier élément.
 * @param <U> Le type du deuxième élément.
 * @param <V> Le type du troisième élément.
 *
 */
public class Triplet<T, U, V> {

  private final T first;
  private final U second;
  private final V third;

  /**
   * @param first Le premier élément du triplet.
   * @param second Le deuxième élément du triplet.
   * @param third Le troisième élément du triplet.
   */
  public Triplet(T first, U second, V third) {
    this.first = first;
    this.second = second;
    this.third = third;
  }

  /**
   * @return Le premier élément du triplet.
   */
  public T getFirst() {
    return first;
  }

  /**
   * @return Le deuxième élément du triplet.
   */
  public U getSecond() {
    return second;
  }

  /**
   * @return Le troisième élément du triplet.
   */
  public V getThird() {
    return third;
  }



}
