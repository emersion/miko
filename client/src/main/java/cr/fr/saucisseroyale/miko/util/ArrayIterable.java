package cr.fr.saucisseroyale.miko.util;

import java.util.Iterator;
import java.util.NoSuchElementException;

/**
 * Un iterable pour simplifier l'itération sur un tableau.
 *
 * @param <T> Le type de données sur lequel itérer.
 *
 */
public class ArrayIterable<T> implements Iterable<T> {

  private static class ArrayIterator<T> implements Iterator<T> {

    private final T[] data;
    private int position = 0;

    public ArrayIterator(T[] data) {
      this.data = data;
    }

    @Override
    public boolean hasNext() {
      return position < data.length;
    }

    @Override
    public T next() {
      if (!hasNext()) {
        throw new NoSuchElementException("End of iterable reached");
      }
      return data[position++];
    }

  }

  private Iterator<T> iterator;

  public ArrayIterable(T[] data) {
    iterator = new ArrayIterator<>(data);
  }

  @Override
  public Iterator<T> iterator() {
    return iterator;
  }

}
