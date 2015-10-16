package cr.fr.saucisseroyale.miko.protocol;

import cr.fr.saucisseroyale.miko.util.IdSaver;
import cr.fr.saucisseroyale.miko.util.UniquelyIdentifiable;


/**
 * Un type de message (entrant ou sortant).
 *
 */
public enum MessageType implements UniquelyIdentifiable {

  // @noformatting
  PING(0), PONG(1),
  EXIT(2),
  LOGIN(3), LOGIN_RESPONSE(4),
  REGISTER(5), REGISTER_RESPONSE(6),
  META_ACTION(7),
  TERRAIN_UPDATE(8), TERRAIN_REQUEST(9),
  ENTITIES_UPDATE(10), ENTITY_UPDATE(11),
  ACTIONS(12), ACTION(13),
  ENTITY_CREATE(14), ENTITY_DESTROY(15),
  CHAT_SEND(16), CHAT_RECEIVE(17),
  VERSION(18);
  // @formatting

  static {
    IdSaver.register(MessageType.class);
  }

  private final int id;

  private MessageType(int id) {
    assert id < 1 << 8 && id >= 0 : "l'identifiant de l'enum est trop petit ou trop grand";
    this.id = id;
  }

  /**
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
