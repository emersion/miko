package cr.fr.saucisseroyale.miko.util;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

/**
 * Utilitaire permettant de retrouver des enums par leur id.
 * <p>
 * Utilisation : enregistrer la classe avec {@link #register(Class, int)} puis récupérer des
 * éléments avec {@link #getValue(Class, int)}.
 */
public class IdSaver {

  // Map (T) : Class<T extends Enum<T>> -> T[]
  private static Map<Class<?>, Object[]> map = new HashMap<>();

  // Classe statique
  private IdSaver() {
    throw new IllegalArgumentException("This class cannot be instantiated");
  }

  /**
   * Renvoit l'id correspondant à un élément d'une enum précédemment sauvegardée, ou null s'il n'y a
   * pas d'élément correspondant.
   *
   * @param enumeration L'énumération identifiable dont on souhaite retrouver un élément.
   * @param id La valeur de l'élément à retrouver.
   * @return L'élément correspondant à id pour cette classe.
   * @throws IllegalArgumentException Si la classe n'a pas été sauvegardée avant utilisation.
   * @throws IndexOutOfBoundsException Si l'id est négatif ou trop grand.
   */
  public static <T extends Enum<T> & UniquelyIdentifiable> T getValue(Class<T> enumeration, int id) {
    Object[] ordered = map.get(enumeration);
    if (ordered == null) {
      throw new IllegalArgumentException("Enum " + enumeration.getCanonicalName()
          + " has not been registered before use");
    }
    if (id < 0 || id > ordered.length) {
      throw new IndexOutOfBoundsException("Id " + id + " must be between 0 and "
          + (ordered.length - 1) + " (inclusive)");
    }
    // On sait que l'élément récupéré est du type de sa clef
    @SuppressWarnings("unchecked")
    T value = (T) ordered[id];
    // value == null <=> pas d'élément trouvé
    return value;
  }

  /**
   * Sauvegarde les enums pour retrouver par id ultérieurement.
   *
   * @param enumeration L'énumération identifiable à sauvegarder.
   * @param maxLength La taille maximum des id. Doit être positive et toujours supérieure ou égale
   *        aux identifiants.
   * @throws IllegalArgumentException Si les identifiants ou maxLength ne respectent pas les
   *         conditions susmentionnées.
   */
  public static <T extends Enum<T> & UniquelyIdentifiable> void register(Class<T> enumeration,
      int maxLength) {
    if (maxLength < 0) {
      throw new IllegalArgumentException("Maxlength must be positive");
    }
    Object[] ordered = new Object[maxLength];
    Arrays.fill(ordered, null);
    T[] constants = enumeration.getEnumConstants();
    for (T constant : constants) {
      int id = constant.getId();
      if (id < 0 || id >= ordered.length) {
        throw new IllegalArgumentException("Id " + id + " must be between 0 and "
            + (ordered.length - 1) + " (inclusive)");
      }
      ordered[id] = constant;
    }
    map.put(enumeration, ordered);
  }
}
