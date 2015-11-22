package message

// An action unique identifier.
type ActionId uint16

// An action. Contains its id and the id of the entity that triggered it.
type Action struct {
	Id        ActionId
	Initiator EntityId
	Params    []interface{}
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
