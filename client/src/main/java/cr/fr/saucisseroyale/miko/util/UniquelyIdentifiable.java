package cr.fr.saucisseroyale.miko.util;

/**
 * Un objet pouvant être identifié par un entier positif et unique. <b>Tous les éléments
 * implémentant cette interface doivent renvoyer une valeur différente.</b>
 *
 */
public interface UniquelyIdentifiable {

  /**
   * Renvoie l'identifiant unique associé à cet élément.
   *
   * @return L'entier positif et unique représentant l'identifiant associé à cet élément.
   */
  public int getId();

}
