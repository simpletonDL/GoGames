package protocol

var BodyKind = struct {
	Box  uint8
	Hero uint8
}{
	Box:  0,
	Hero: 1,
}

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
	ImageKind     uint8
	Width, Height float64
}
