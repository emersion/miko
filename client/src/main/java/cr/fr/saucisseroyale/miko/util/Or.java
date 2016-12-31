package cr.fr.saucisseroyale.miko.util;

/**
 * Un élément pouvant contenir un élément d'une classe, OU (exclusivement) d'une autre classe.
 *
 * @param <T> La première classe possible pour cet élément.
 * @param <U> La deuxième classe possible pour cet élément
 */
@SuppressWarnings("unchecked")
public class Or<T, U> {
  private final Object element;
  private final boolean first;

  /**
   * @param element L'élément à mettre dans cet objet.
   */
  public Or(Object element, boolean first) {
    this.element = element;
    this.first = first;
  }

  /**
   * @return L'élement stocké dans cet objet, ou null s'il était du type de la deuxième classe.
   */
  public T getAsFirst() {
    return first ? (T) element : null;
  }

  /**
   * @return L'élement stocké dans cet objet, ou null s'il était du type de la première classe.
   */
  public U getAsSecond() {
    return first ? null : (U) element;
  }

  /**
   * @return true si l'élément est du type de la première classe, false sinon.
   */
  public boolean isFirst() {
    return first;
  }
}
