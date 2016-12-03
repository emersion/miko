package cr.fr.saucisseroyale.miko.util;

import java.util.Iterator;
import java.util.NoSuchElementException;

/**
 * Un iterable pour simplifier l'itération sur un seul élément.
 *
 * @param <T> Le type de la donnée sur laquelle itérer.
 */
public class SingleElementIterable<T> implements Iterable<T> {
  private Iterator<T> iterator;

  public SingleElementIterable(T element) {
    iterator = new SingleElementIterator<>(element);
  }

  @Override
  public Iterator<T> iterator() {
    return iterator;
  }

  private static class SingleElementIterator<T> implements Iterator<T> {
    private final T element;
    private boolean returnedElement = false;

    public SingleElementIterator(T element) {
      this.element = element;
    }

    @Override
    public boolean hasNext() {
      return !returnedElement;
    }

    @Override
    public T next() {
      if (!hasNext()) {
        throw new NoSuchElementException("End of iterable reached");
      }
      return element;
    }
  }
}
