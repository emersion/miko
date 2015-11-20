package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// A delta (not to be confused with @delthas) stores all data about an entity
// change: not only its diff if it has been updated, but also the time when the
// change has been made and the entity itself.
type Delta struct {
	tick      message.AbsoluteTick
	requested bool
	EntityId  message.EntityId

	Diff *message.EntityDiff
	From *Entity
	To   *Entity
}

// Get the tick when this delta was executed.
func (d *Delta) GetTick() message.AbsoluteTick {
	return d.tick
}

// Check that this delta has been requested by a client, rather than generated
// by the server.
func (d *Delta) Requested() bool {
	return d.requested
}

// Get a request triggering this delta's changes.
func (d *Delta) Request() Request {
	req := newRequest(d.tick)
	req.requested = d.requested

	// This delta has been created when a request was accepted. Set accepted to
	// true to prevent it to create another delta.
	req.accepted = true

	if d.From != nil && d.To != nil { // Update
		if d.Diff == nil {
			panic("Cannot build update request from delta: no diff available")
		}

		return &UpdateRequest{req, d.To, d.Diff}
	} else if d.To != nil { // Create
		return &CreateRequest{req, d.To}
	} else if d.From != nil { // Delete
		return &DeleteRequest{req, d.From.Id}
	}

	panic("Cannot build request from delta: from and to are empty")
}

// Get the inverse of this delta.
func (d *Delta) Inverse() *Delta {
	return &Delta{
		tick:      d.tick,
		requested: d.requested,
		EntityId:  d.EntityId,
		Diff:      d.Diff,
		From:      d.To,
		To:        d.From,
	}
}

// Ensure that there is only one delta per entity in the deltas list.
// Warning: some information is lost, for instance if an entity is created then
// deleted, an empty delta is added to the list.
func flattenDeltas(deltas []*Delta) []*Delta {
	m := make(map[message.EntityId]*Delta)

	for _, d := range deltas {
		if current, ok := m[d.EntityId]; ok {
			// TODO: check that current.To and d.From are compatible
			current.To = d.To

			if d.Diff != nil {
				current.Diff.Merge(d.Diff)
			}
		} else {
			m[d.EntityId] = &Delta{
				EntityId: d.EntityId,
				Diff:     message.NewEntityDiff(),
				From:     d.From,
				To:       d.To,
			}

			if d.Diff != nil {
				m[d.EntityId].Diff.Merge(d.Diff)
			}
		}
	}

	l := make([]*Delta, len(m))
	i := 0
	for _, d := range m {
		l[i] = d
		i++
	}

	return l
}

// Convert a deltas list to a diff pool. Assumes that the deltas list has been
// flattened.
func deltasToDiffPool(deltas []*Delta) *message.EntityDiffPool {
	pool := message.NewEntityDiffPool()

	for _, d := range deltas {
		if d.From != nil && d.To != nil { // Update
			from := d.From.ToMessage()
			to := d.To.ToMessage()

			if from.EqualsWithDiff(to, d.Diff) {
				continue
			}

			pool.Updated[d.To.ToMessage()] = d.Diff
		} else if d.To != nil { // Create
			pool.Created = append(pool.Created, d.To.ToMessage())
		} else if d.From != nil { // Delete
			pool.Deleted = append(pool.Deleted, d.EntityId)
		}
	}

	return pool
}
