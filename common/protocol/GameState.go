package protocol

type BodyKind uint8

type DirectionKind bool

const (
	BodyKindBox = BodyKind(iota)
	BodyKindHero
	BodyKindBullet
	BodyKindPlatform
)

const (
	DirectionKindRight = DirectionKind(true)
	DirectionKindLeft  = DirectionKind(false)
)

type GameState struct {
	Objects []GameObject
}

func NewEmptyGameState() GameState {
	return GameState{
		Objects: []GameObject{},
	}
}

type GameObject struct {
	XPos          float64
	YPos          float64
	Angel         float64
	BodyKind      BodyKind
	Width, Height float64
	Direction     DirectionKind
}
