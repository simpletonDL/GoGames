package client

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
)

func MakeImageOptions(image *ebiten.Image, width float64, height float64, xPos float64, yPos float64, angel float64, inverseX bool) *ebiten.DrawImageOptions {
	originWidth, originHeight := image.Bounds().Dx(), image.Bounds().Dy()
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(width/float64(originWidth), height/float64(originHeight))
	options.GeoM.Translate(-width/2, -height/2)
	if inverseX {
		options.GeoM.Scale(-1, 1)
	}
	options.GeoM.Rotate(angel)
	options.GeoM.Translate(xPos, yPos)
	return options
}

func getWeaponImage(kind protocol.WeaponKind) *ebiten.Image {
	switch kind {
	case protocol.WeaponKindDefault:
		return DefaultWeaponImage
	case protocol.WeaponKindSniperRifle:
		return SniperRifleWeaponImage
	case protocol.WeaponKindMachineGun:
		return MachineGunWeaponImage
	default:
		panic(fmt.Sprintf("Unknown weapon kind: %d", kind))
	}
}

func Render(image *ebiten.Image, state *protocol.GameState) {
	scaleX := float64(image.Bounds().Dx()) / settings.WorldWidth
	scaleY := float64(image.Bounds().Dy()) / settings.WorldHeight

	for _, obj := range state.Objects {
		var weaponImage *ebiten.Image = nil

		var objImage *ebiten.Image
		switch obj.BodyKind {
		case protocol.BodyKindBox:
			objImage = BoxImage
		case protocol.BodyKindHero:
			objImage = HeroImage
			weaponImage = getWeaponImage(obj.WeaponKind)
		case protocol.BodyKindBullet:
			objImage = BulletImage
		case protocol.BodyKindPlatform:
			objImage = PlatformImage
		default:
			panic(fmt.Sprintf("Unknown object kind %d", obj.BodyKind))
		}

		width := obj.Width * scaleX
		height := obj.Height * scaleY
		xPos := obj.XPos * scaleX
		yPos := float64(image.Bounds().Dy()) - obj.YPos*scaleY
		angel := -obj.Angel
		inverseX := obj.Direction == protocol.DirectionKindLeft
		objOptions := MakeImageOptions(objImage, width, height, xPos, yPos, angel, inverseX)
		image.DrawImage(objImage, objOptions)

		if weaponImage != nil {
			// weapon image height should be hero height
			scale := float64(height) / float64(weaponImage.Bounds().Dy())
			weaponWidthX := float64(weaponImage.Bounds().Dx()) * scale

			// TODO: adjust weapon by angel
			var weaponAdjustmentX float64
			if inverseX {
				weaponAdjustmentX = -width / 2
			} else {
				weaponAdjustmentX = width / 2
			}
			weaponOptions := MakeImageOptions(weaponImage, weaponWidthX, height, xPos+weaponAdjustmentX, yPos, angel, inverseX)
			image.DrawImage(weaponImage, weaponOptions)
		}
	}
}
