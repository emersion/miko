package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;


/**
 * Représente un type de message (entrant ou sortant) sous forme de byte.
 *
 */
public enum MessageType implements UniquelyIdentifiable {

  // @noformatting
  PING(0), PONG(1),
  EXIT(2),
  LOGIN(3), LOGIN_RESPONSE(4),
  REGISTER(5), REGISTER_RESPONSE(6),
  METAACTION(7),
  TERRAINUPDATE(8), TERRAINREQUEST(9),
  ENTITIESUPDATE(10), ENTITYUPDATE(11),
  ACTIONS(12), ACTION(13),
  ENTITIESCREATE(14), ENTITIESDESTROY(15),
  CHATSEND(16), CHATRECEIVE(17);
  // @formatting

  static {
    IdSaver.register(MessageType.class, 1 << 8);
  }

  private final int id;

  private MessageType(int id) {
    assert id < 1 << 8 && id >= 0;
    this.id = id;
  }

  /**
   * Renvoit l'identifiant correspondant à un message.
   *
   * @return L'identifiant correspondant au type du message.
   */
  @Override
  public int getId() {
    return id;
  }

  /**
   * Renvoit le type d'un message correspondant à un identifiant.
   *
   * @param id L'identifiant correspondant au type du message.
   * @return Le type du message, ou null s'il n'y a pas de message correspondant.
   */
  public static MessageType getType(int id) {
    return IdSaver.getValue(MessageType.class, id);
  }
}
