package cr.fr.saucisseroyale.miko;


/**
 * Représente un type de paquet (entrant ou sortant) sous forme de byte.
 *
 */
enum PacketType {

  // On ne stocke pas les cas de réponse (+ 2^7), ils correspondent à leur version normale donc pas
  // besoin de les faire en double. On fera la différence entre les deux lors du parsing.
  // @noformatting
  PING(0), EXIT(1),
  LOGIN(2), REGISTER(3), METAACTION(4),
  TERRAINUPDATE(5), TERRAINREQUEST(6),
  ENTITIESUPDATE(7), ENTITYUPDATE(8),
  ACTIONS(9), ACTION(10),
  ENTITIESCREATE(11), ENTITIESDESTROY(12),
  CHATSEND(13), CHATRECEIVE(14);
  // @formatting

  private final byte id;

  private PacketType(int id) {
    assert (id < (1 << 7));
    this.id = (byte) id;
  }

  /**
   * @return Le byte correspondant au type du paquet.
   */
  public byte getId() {
    return id;
  }
}
