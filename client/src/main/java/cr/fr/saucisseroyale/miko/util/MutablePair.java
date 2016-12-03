package cr.fr.saucisseroyale.miko.util;

/**
 * Une paire générique <b>mutable</b> de deux éléments (2-tuple).
 *
 * @param <T> Le type du premier élément.
 * @param <U> Le type du deuxième élément.
 */
public class MutablePair<T, U> {
  private T first;
  private U second;

  /**
   * @param first  Le premier élément de la paire.
   * @param second Le deuxième élément de la paire.
   */
  public MutablePair(T first, U second) {
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
   * @param first Le premier élément de la paire à définir.
   */
  public void setFirst(T first) {
    this.first = first;
  }

  /**
   * @return Le deuxième élément de la paire.
   */
  public U getSecond() {
    return second;
  }

  /**
   * @param second Le deuxième élément de la paire à définir.
   */
  public void setSecond(U second) {
    this.second = second;
  }

  /**
   * Une MutablePair avec long comme type de premier élément.
   * <p>
   * Sert à éviter l'autoboxing de la primitive long.
   *
   * @param <U> Le type du second élément.
   * @see MutablePair
   */
  public static class Long<U> {
    private long first;
    private U second;

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
     * @param first Le premier élément de la paire à définir.
     */
    public void setFirst(long first) {
      this.first = first;
    }

    /**
     * @return Le deuxième élément de la paire.
     */
    public U getSecond() {
      return second;
    }

    /**
     * @param second Le deuxième élément de la paire à définir.
     */
    public void setSecond(U second) {
      this.second = second;
    }
  }
}
