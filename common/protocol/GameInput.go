package protocol

var InputCommandKind = struct {
	MouseClick   uint8
	MoveHero     uint8
	MakeShoot    uint8
	RotateHero   uint8
	ChangeWeapon uint8
}{
	MouseClick:   0,
	MoveHero:     1,
	MakeShoot:    2,
	RotateHero:   3,
	ChangeWeapon: 4,
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

func NewMakeShootCommand() ClientInputCommand {
	return ClientInputCommand{
		Id:        InputCommandKind.MakeShoot,
		FloatArgs: nil,
		IntArgs:   nil,
	}
}

var RotateHeroKind = struct {
	Left  int // counterclockwise
	Right int // clockwise
}{
	Left:  0,
	Right: 1,
}

func NewRotateHeroCommand(rotateKind int) ClientInputCommand {
	return ClientInputCommand{
		Id: InputCommandKind.RotateHero,
		IntArgs: map[string]int{
			"kind": rotateKind,
		},
	}
}

// For debug purposes
func NewChangeWeaponCommand(kind WeaponKind) ClientInputCommand {
	return ClientInputCommand{
		Id: InputCommandKind.ChangeWeapon,
		IntArgs: map[string]int{
			"kind": int(kind),
		},
	}
}
