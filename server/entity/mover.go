package entity

import(
	"time"

	"git.emersion.fr/saucisse-royale/miko/server/message"
)

// A service that moves entities
type Mover struct {
	lastUpdates map[message.EntityId]int64
}

// Compute an entity's new position
func (m *Mover) UpdateEntity(entity *message.Entity) *message.Position {
	last := m.lastUpdates[entity.Id]
	now := time.Now().UnixNano()
	m.lastUpdates[entity.Id] = now
	dt := float64(now - last) / 1000 / 1000 // Convert to seconds
	return entity.Speed.GetNextPosition(entity.Position, dt)
}
