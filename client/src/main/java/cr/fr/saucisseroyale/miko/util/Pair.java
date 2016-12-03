package cr.fr.saucisseroyale.miko.util;

/**
 * Une paire générique immutable de deux éléments (2-tuple).
 *
 * @param <T> Le type du premier élément.
 * @param <U> Le type du deuxième élément.
 */
public class Pair<T, U> {
  private final T first;
  private final U second;

  /**
   * @param first  Le premier élément de la paire.
   * @param second Le deuxième élément de la paire.
   */
  public Pair(T first, U second) {
    this.first = first;
    this.second = second;
  }

  /**
   * @return Le premier élément de la paire.
   */
  public T getFirst() {
    return first;
  }

  /**
   * @return Le deuxième élément de la paire.
   */
  public U getSecond() {
    return second;
  }

  /**
   * Une Pair avec int comme type de premier élément.
   * <p>
   * Sert à éviter l'autoboxing de la primitive int.
   *
   * @param <U> Le type du second élément.
   * @see Pair
   */
  public static class Int<U> {
    private final int first;
    private final U second;

    /**
     * @param first  Le premier élément de la paire.
     * @param second Le deuxième élément de la paire.
     */
    public Int(int first, U second) {
      this.first = first;
      this.second = second;
    }

    /**
     * @return Le premier élément de la paire.
     */
    public int getFirst() {
      return first;
    }

    /**
     * @return Le deuxième élément de la paire.
     */
    public U getSecond() {
      return second;
    }
  }

  /**
   * Une Pair avec long comme type de premier élément.
   * <p>
   * Sert à éviter l'autoboxing de la primitive long.
   *
   * @param <U> Le type du second élément.
   * @see Pair
   */
  public static class Long<U> {
    private final long first;
    private final U second;

    /**
     * @param first  Le premier élément de la paire.
     * @param second Le deuxième élément de la paire.
     */
    public Long(long first, U second) {
      this.first = first;
      this.second = second;
    }

    /**
     * @return Le premier élément de la paire.
     */
    public long getFirst() {
      return first;
    }

    /**
     * @return Le deuxième élément de la paire.
     */
    public U getSecond() {
      return second;
    }
  }

  /**
   * Une Pair avec float comme type pour les deux éléments.
   * <p>
   * Sert à éviter l'autoboxing des primitivess float.
   *
   * @see Pair
   */
  public static class DoubleFloat {
    private final float first;
    private final float second;

    /**
     * @param first  Le premier élément de la paire.
     * @param second Le deuxième élément de la paire.
     */
    public DoubleFloat(float first, float second) {
      this.first = first;
      this.second = second;
    }

    /**
     * @return Le premier élément de la paire.
     */
    public float getFirst() {
      return first;
    }

    /**
     * @return Le deuxième élément de la paire.
     */
    public float getSecond() {
      return second;
    }
  }
}
