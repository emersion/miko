package cr.fr.saucisseroyale.miko.util;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

/**
 * Utilitaire permettant de retrouver des enums par leur id.
 * <p>
 * Utilisation : enregistrer la classe avec {@link #register(Class)} puis récupérer des éléments
 * avec {@link #getValue(Class, int)}.
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
   */
  public static <T extends Enum<T> & UniquelyIdentifiable> T getValue(Class<T> enumeration, int id) {
    Object[] ordered = map.get(enumeration);
    if (ordered == null) {
      throw new IllegalArgumentException("Enum " + enumeration.getCanonicalName()
          + " has not been registered before use");
    }
    if (id < 0 || id > ordered.length) {
      return null;
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
   */
  public static <T extends Enum<T> & UniquelyIdentifiable> void register(Class<T> enumeration) {
    T[] constants = enumeration.getEnumConstants();
    int maxId = -1;
    for (T constant : constants) {
      int id = constant.getId();
      if (id < 0) {
        throw new IllegalArgumentException("Id " + id + " must positive");
      }
      if (id > maxId)
        maxId = id;
    }
    Object[] ordered = new Object[maxId + 1];
    Arrays.fill(ordered, null);
    for (T constant : constants) {
      int id = constant.getId();
      ordered[id] = constant;
    }
    map.put(enumeration, ordered);
  }
}
