package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.List;


class EngineMessage implements Comparable<EngineMessage> {

  enum Type {
    PLAYER_JOINED(1), PLAYER_LEFT(2), CHUNK_UPDATE(3), ACTIONS_DONE(4), ENTITY_DESTROY(5), ENTITY_CREATE(6), ENTITIES_UPDATE(7);

    private final int order;

    private Type(int order) {
      this.order = order;
    }
  }

  private final long tick;
  private final Type type;
  private final Object[] data;

  private EngineMessage(long tick, Type type, Object[] data) {
    this.tick = tick;
    this.type = type;
    this.data = data;
  }

  public static EngineMessage newPlayerJoinedMessage(long tick, int playerId, String name) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[] {playerId, name});
  }

  public static EngineMessage newPlayerLeftMessage(long tick, int playerId) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[] {playerId});
  }

  public static EngineMessage newChunkUpdateMessage(long tick, ChunkPoint chunkPoint, Chunk chunk) {
    return new EngineMessage(tick, Type.CHUNK_UPDATE, new Object[] {chunkPoint, chunk});
  }

  public static EngineMessage newActionsMessage(long tick, List<Pair<Integer, Action>> actions) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[] {actions});
  }

  public static EngineMessage newEntityDestroyMessage(long tick, int entityId) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[] {entityId});
  }

  public static EngineMessage newEntityCreateMessage(long tick, EntityDataUpdate entityDataUpdate) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[] {entityDataUpdate});
  }

  public static EngineMessage newEntitiesUpdateMessage(long tick, List<EntityDataUpdate> entitiesUpdateList) {
    return new EngineMessage(tick, Type.PLAYER_JOINED, new Object[] {entitiesUpdateList});
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

}
