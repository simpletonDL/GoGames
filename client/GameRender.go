package client

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
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
	case protocol.WeaponKindCarbine:
		return CarbineWeaponImage
	default:
		panic(fmt.Sprintf("Unknown weapon kind: %d", kind))
	}
}

func Render(image *ebiten.Image, state *protocol.GameState) {
	scaleX := float64(image.Bounds().Dx()) / settings.WorldWidth
	scaleY := float64(image.Bounds().Dy()) / settings.WorldHeight

	for _, obj := range state.Objects {
		var weaponImage *ebiten.Image = nil
		var nickname string

		var objImage *ebiten.Image
		switch obj.BodyKind {
		case protocol.BodyKindBox:
			objImage = BoxImage
		case protocol.BodyKindHero:
			objImage = HeroImage
			weaponImage = getWeaponImage(obj.WeaponKind)
			nickname = obj.Nickname
		case protocol.BodyKindBullet:
			objImage = BulletImage
		case protocol.BodyKindPlatform:
			objImage = PlatformImage
		case protocol.BodyKindWeaponBox:
			objImage = WeaponBoxImage
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

		if nickname != "" {
			face := &text.GoTextFace{Source: MainFont, Size: 8}
			w, h := text.Measure(nickname, face, 0)
			opts := text.DrawOptions{}
			opts.GeoM.Translate(xPos-w/2, yPos-(height)/2-h)
			if obj.Team == protocol.BlueTeam {
				opts.ColorScale.Scale(0, 191, 255, 1)
			} else {
				opts.ColorScale.Scale(255, 0, 0, 1)
			}
			text.Draw(image, nickname, face, &opts)
		}
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

	// Add text info about hero: lifes, weapon bullets, e.t.c.
	players := utils.Filter(state.Objects, func(object protocol.GameObject) bool {
		return object.BodyKind == protocol.BodyKindHero
	})
	blueTeamPlayers := utils.Filter(players, func(object protocol.GameObject) bool { return object.Team == protocol.BlueTeam })
	redTeamPlayers := utils.Filter(players, func(object protocol.GameObject) bool { return object.Team == protocol.RedTeam })
	addPlayerTextInfo(image, blueTeamPlayers, protocol.BlueTeam)
	addPlayerTextInfo(image, redTeamPlayers, protocol.RedTeam)
}

func addPlayerTextInfo(image *ebiten.Image, players []protocol.GameObject, team protocol.TeamKind) {
	for i, player := range players {
		var weaponAvailableBullets string
		if player.WeaponAvailableBullets > 100_000 {
			weaponAvailableBullets = "inf"
		} else {
			weaponAvailableBullets = fmt.Sprintf("%d", player.WeaponAvailableBullets)
		}

		playerTextInfo := fmt.Sprintf("%s L=%d B=%d/%d (%s)",
			utils.AdjustString(player.Nickname, 15),
			player.LivesCount,
			player.WeaponAvailableBulletsInMagazine,
			player.WeaponMagazineCapacity,
			weaponAvailableBullets,
		)
		// TODO: place text in respect of team (after introducing teams)
		yPos := i * 25
		xPos := 0
		switch team {
		case protocol.BlueTeam:
			xPos = 0
		case protocol.RedTeam:
			xPos = settings.ScreenWidth / 2
		}

		opts := text.DrawOptions{}
		opts.GeoM.Translate(float64(xPos), float64(yPos))
		text.Draw(image, playerTextInfo, &text.GoTextFace{Source: MainFont, Size: 15}, &opts)
	}
}
