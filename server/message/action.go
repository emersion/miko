package message

// An action unique identifier.
type ActionId uint16

// An action. Contains its id and the id of the entity that triggered it.
type Action struct {
	Id        ActionId
	Initiator EntityId
	Params    []interface{}
}

// Check if this action is the same as another one.
func (a *Action) Equals(other *Action) bool {
	if a.Id != other.Id || a.Initiator != other.Initiator {
		return false
	}
	if len(a.Params) != len(other.Params) {
		return false
	}
	for i, param := range a.Params {
		otherParam := other.Params[i]
		if param != otherParam {
			return false
		}
	}
	return true
}

func NewAction() *Action {
	return &Action{}
}

// An action service.
type ActionService interface {
	// Execute an action.
	Execute(a *Action, t AbsoluteTick) Request

	// Check if some actions have been executed.
	IsDirty() bool

	// Flush the actions history.Return the current one and replace it by a new
	// empty one.
	Flush() []*Action
}
