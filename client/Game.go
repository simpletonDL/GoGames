package client

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"log"
	"net"
)

type Game struct {
	Conn      net.Conn
	GameState protocol.GameState
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		coordinatesMapper := engine.NewCoordinatesMapper(settings.WorldWidth, settings.WorldHeight,
			settings.ScreenWidth, settings.ScreenHeight)

		worldX, worldY := coordinatesMapper.ScreenToWorld(float64(x), float64(y))
		fmt.Printf("%f %f\n", worldX, worldY)
		message := protocol.NewMouseClick(worldX, worldY)

		// Serialize the message to JSON
		encoder := json.NewEncoder(g.Conn)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v", err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	backgroundOptions := MakeImageOptions(BackgroundImage, settings.ScreenWidth, settings.ScreenHeight,
		settings.ScreenWidth/2, settings.ScreenHeight/2, 0)
	screen.DrawImage(BackgroundImage, backgroundOptions)
	Render(screen, &g.GameState)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return settings.ScreenWidth, settings.ScreenHeight
}
