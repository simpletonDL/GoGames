package protocol

var InputCommandKind = struct {
	MouseClick uint8
	MoveHero   uint8
}{
	MouseClick: 0,
	MoveHero:   1,
}

type GameInput []ClientInputCommand

type ClientInputCommand struct {
	Id        uint8
	FloatArgs map[string]float64
	IntArgs   map[string]int
}

func NewMouseClickCommand(worldX float64, worldY float64) ClientInputCommand {
	return ClientInputCommand{
		Id: InputCommandKind.MouseClick,
		FloatArgs: map[string]float64{
			"x": worldX,
			"y": worldY,
		},
	}
}

var MoveHeroKind = struct {
	Left  int
	Right int
	Up    int
	Down  int
}{
	Left:  0,
	Right: 1,
	Up:    2,
	Down:  3,
}

func NewMoveHeroCommand(moveKind int) ClientInputCommand {
	return ClientInputCommand{
		Id: InputCommandKind.MoveHero,
		IntArgs: map[string]int{
			"kind": moveKind,
		},
	}
}
