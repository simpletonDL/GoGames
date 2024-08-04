package client

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
)

func MakeImageOptions(image *ebiten.Image, width float64, height float64, xPos float64, yPos float64, angel float64) *ebiten.DrawImageOptions {
	originWidth, originHeight := image.Bounds().Dx(), image.Bounds().Dy()
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(width/float64(originWidth), height/float64(originHeight))
	options.GeoM.Translate(-width/2, -height/2)
	options.GeoM.Rotate(angel)
	options.GeoM.Translate(xPos, yPos)
	return options
}

func Render(image *ebiten.Image, state *protocol.GameState) {
	scaleX := float64(image.Bounds().Dx()) / settings.WorldWidth
	scaleY := float64(image.Bounds().Dy()) / settings.WorldHeight

	for _, obj := range state.Objects {
		var objImage *ebiten.Image
		switch obj.ImageKind {
		case protocol.BodyKind.Box:
			objImage = BoxImage
		case protocol.BodyKind.Hero:
			objImage = HeroImage
		case protocol.BodyKind.Bullet:
			objImage = BulletImage
		default:
			panic(fmt.Sprintf("Unknown object kind %d", obj.ImageKind))
		}

		objOptions := MakeImageOptions(objImage, obj.Width*scaleX, obj.Height*scaleY,
			obj.XPos*scaleX, float64(image.Bounds().Dy())-obj.YPos*scaleY, -obj.Angel)
		image.DrawImage(objImage, objOptions)
	}
}
