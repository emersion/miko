package cr.fr.saucisseroyale.miko.network;


/**
 * Représente un type de message (entrant ou sortant) sous forme de byte.
 *
 */
enum MessageType {

  // On ne stocke pas les cas de réponse (+ 2^7), ils correspondent à leur version normale donc pas
  // besoin de les faire en double. On fera la différence entre les deux lors du parsing.
  // TODO maj care le protocole a changé
  // @noformatting
  PING(0), EXIT(1),
  LOGIN(2), REGISTER(3), METAACTION(4),
  TERRAINUPDATE(5), TERRAINREQUEST(6),
  ENTITIESUPDATE(7), ENTITYUPDATE(8),
  ACTIONS(9), ACTION(10),
  ENTITIESCREATE(11), ENTITIESDESTROY(12),
  CHATSEND(13), CHATRECEIVE(14);
  // @formatting

  private static MessageType[] types = new MessageType[1 << 7];
  private final int id;

  static {
    for (MessageType type : MessageType.values()) {
      types[type.id] = type;
    }
  }

  private MessageType(int id) {
    assert id < 1 << 7 && id > 0;
    this.id = id;
  }

  /**
   * Renvoit le nombre correspondant à un message.
   *
   * @param isResponse Si le type du message sera pour une réponse ou pour un simple envoi.
   * @return Le nombre correspondant au type du message.
   */
  public int getId(boolean isResponse) {
    return id;
  }

  /**
   * Renvoit le type d'un message correspondant à un nombre.
   *
   * @param id Le nombre correspondant au type du message.
   * @return Le type du message, ou null s'il n'y a pas de message correspondant.
   */
  public static MessageType getType(int id) {
    if (id >= 1 << 8 || id < 0) {
      throw new IllegalArgumentException("Id has to be less than " + (1 << 8));
    }
    return types[id & 0x7F]; // On réduit le cas de réponse au cas normal (-2^7 si besoin)
  }

}
