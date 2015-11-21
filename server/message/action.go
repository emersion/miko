package message

// An action unique identifier.
type ActionId uint16

// An action. Contains its id and the id of the entity that triggered it.
type Action struct {
	Id        ActionId
	Initiator EntityId
}

// An action service.
type ActionService interface{}
