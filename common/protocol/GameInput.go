package protocol

const (
	MouseClickID = iota
)

type GameInput []InputCommand

type InputCommand struct {
	Id   int
	Args map[string]float64
}

func NewMouseClick(worldX float64, worldY float64) InputCommand {
	return InputCommand{
		Id: MouseClickID,
		Args: map[string]float64{
			"x": worldX,
			"y": worldY,
		},
	}
}
