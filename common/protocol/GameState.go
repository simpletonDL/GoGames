package protocol

const (
	BoxKind = iota
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
	ImageKind     float64
	Width, Height float64
}
