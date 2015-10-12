package message

type EntityId uint16

type Position struct {
	BX BlockCoord
	BY BlockCoord
	X PointCoord
	Y PointCoord
}

type Speed struct {
	Angle float32
	Norm float32
}

type Entity struct {
	Id EntityId
	Position *Position
	Speed *Speed
}

type EntityService interface {}
