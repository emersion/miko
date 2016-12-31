package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.util.Pair;
import cr.fr.saucisseroyale.miko.util.Pair.Int;
import edu.umd.cs.findbugs.annotations.SuppressFBWarnings;

import java.util.List;

/**
 * Un message adressé à l'engine du jeu, stocké sous forme d'objet pour délayer son exécution.
 * <p>
 * Note: L'ordre naturel donné par {@link #compareTo(EngineMessage)} n'est pas cohérent avec
 * {@link #equals(Object)} et ne sert que pour comparer l'ordre dans lequel traiter les messages.
 */
@SuppressFBWarnings("EQ_COMPARETO_USE_OBJECT_EQUALS")
final class EngineMessage implements Comparable<EngineMessage> {
  private final long tick;
  private final Type type;
  private final Object[] data;

  private EngineMessage(long tick, Type type, Object[] data) {
    this.tick = tick;
    this.type = type;
    this.data = data;
  }

  public static EngineMessage newPlayerJoinedMessage(long tick, int playerId, String name) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[]{playerId, name});
  }

  public static EngineMessage newPlayerLeftMessage(long tick, int playerId) {
    return new EngineMessage(tick, Type.PLAYER_LEFT, new Object[]{playerId});
  }

  public static EngineMessage newChunksUpdateMessage(long tick, List<Pair<ChunkPoint, Chunk>> chunks) {
    return new EngineMessage(tick, Type.CHUNKS_UPDATE, new Object[]{chunks});
  }

  public static EngineMessage newActionsMessage(long tick, List<Int<Action>> actions) {
    return new EngineMessage(tick, Type.ACTIONS_DONE, new Object[]{actions});
  }

  public static EngineMessage newEntityDestroyMessage(long tick, int entityId) {
    return new EngineMessage(tick, Type.ENTITY_DESTROY, new Object[]{entityId});
  }

  public static EngineMessage newEntityCreateMessage(long tick, EntityDataUpdate entityDataUpdate) {
    return new EngineMessage(tick, Type.ENTITY_CREATE, new Object[]{entityDataUpdate});
  }

  public static EngineMessage newEntitiesUpdateMessage(long tick, List<EntityDataUpdate> entitiesUpdateList) {
    return new EngineMessage(tick, Type.ENTITIES_UPDATE, new Object[]{entitiesUpdateList});
  }

  public long getTick() {
    return tick;
  }

  public Type getType() {
    return type;
  }

  public Object[] getData() {
    return data;
  }

  @Override
  public int compareTo(EngineMessage o) {
    if (tick < o.tick) {
      return -1;
    }
    if (tick > o.tick) {
      return 1;
    }
    if (type.order < o.type.order) {
      return -1;
    }
    if (type.order > o.type.order) {
      return 1;
    }
    return 0;
  }

  enum Type {
    PLAYER_JOINED(1), PLAYER_LEFT(2), CHUNKS_UPDATE(3), ACTIONS_DONE(4), ENTITY_DESTROY(5), ENTITY_CREATE(6), ENTITIES_UPDATE(7);
    private final int order;

    Type(int order) {
      this.order = order;
    }
  }
}
