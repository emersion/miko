package cr.fr.saucisseroyale.miko.protocol;

/**
 * Un type de paramètres associés à des actions.
 *
 * @see ActionType
 * @see Action
 */
public enum DataType {

  // update InputMessageFactory#readObject, OutputMessageFactory#entityUpdate,
  // OutputMessageFactory#writeAction, et Action en cas de changement
  // @noformatting
  VOID,
  ONE_FLOAT,
  ONE_ENTITY,
  ONE_SHORT,
  ONE_TERRAIN,
  PAIR_FLOAT_ENTITY;
  // @formatting
}
