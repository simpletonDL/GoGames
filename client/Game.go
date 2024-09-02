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
	Encoder   *json.Encoder
	GameState protocol.GameState
}

func NewGame(conn net.Conn) *Game {
	return &Game{
		Conn:      conn,
		Encoder:   json.NewEncoder(conn),
		GameState: protocol.NewEmptyGameState(),
	}
}

func (g *Game) Send(msg any) {
	err := g.Encoder.Encode(msg)
	if err != nil {
		log.Printf("Failed to serialize message to JSON: %v\n", err)
	}
}

func (g *Game) Update() error {
	encoder := json.NewEncoder(g.Conn)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		coordinatesMapper := engine.NewCoordinatesMapper(settings.WorldWidth, settings.WorldHeight,
			settings.ScreenWidth, settings.ScreenHeight)

		worldX, worldY := coordinatesMapper.ScreenToWorld(float64(x), float64(y))
		fmt.Printf("%f %f\n", worldX, worldY)
		message := protocol.NewMouseClickCommand(worldX, worldY)

		// Serialize the message to JSON
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		message := protocol.NewMoveHeroCommand(protocol.MoveHeroKind.Right)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		message := protocol.NewMoveHeroCommand(protocol.MoveHeroKind.Left)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		message := protocol.NewMoveHeroCommand(protocol.MoveHeroKind.Up)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		message := protocol.NewMoveHeroCommand(protocol.MoveHeroKind.Down)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		message := protocol.NewMakeShootCommand()
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyO) {
		message := protocol.NewRotateHeroCommand(protocol.RotateHeroKind.Left)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		message := protocol.NewRotateHeroCommand(protocol.RotateHeroKind.Right)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDigit1) {
		message := protocol.NewChangeWeaponCommand(protocol.WeaponKindDefault)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDigit2) {
		message := protocol.NewChangeWeaponCommand(protocol.WeaponKindSniperRifle)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDigit3) {
		message := protocol.NewChangeWeaponCommand(protocol.WeaponKindMachineGun)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDigit4) {
		message := protocol.NewChangeWeaponCommand(protocol.WeaponKindCarbine)
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		message := protocol.NewReadyToStartCommand()
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		message := protocol.NewNotReadyToStartCommand()
		err := encoder.Encode(message)
		if err != nil {
			log.Printf("Failed to serialize message to JSON: %v\n", err)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	backgroundOptions := MakeImageOptions(BackgroundImage, settings.ScreenWidth, settings.ScreenHeight,
		settings.ScreenWidth/2, settings.ScreenHeight/2, 0, false)
	screen.DrawImage(BackgroundImage, backgroundOptions)
	Render(screen, &g.GameState)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return settings.ScreenWidth, settings.ScreenHeight
}
