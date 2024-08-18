package client

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

var (
	//go:embed assets/stars-background.jpg
	backgroundImageRaw []byte

	//go:embed assets/wooden-box.png
	boxImageRaw []byte

	//go:embed assets/hero.png
	hero []byte

	//go:embed assets/paintball.png
	bullet []byte

	//go:embed assets/platform.png
	platform []byte

	/* Weapons */
	//go:embed assets/weapons/pistol.png
	defaultWeapon []byte

	//go:embed assets/weapons/sniper-rifle.png
	sniperRifleWeapon []byte

	//go:embed assets/weapons/machine-gun.png
	machineGunWeapon []byte
)

var (
	BackgroundImage *ebiten.Image
	BulletImage     *ebiten.Image
	PlatformImage   *ebiten.Image
	BoxImage        *ebiten.Image
	HeroImage       *ebiten.Image

	/* Weapons */

	DefaultWeaponImage     *ebiten.Image
	SniperRifleWeaponImage *ebiten.Image
	MachineGunWeaponImage  *ebiten.Image
)

func LoadImage(bs []byte) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func LoadImages() {
	BackgroundImage = LoadImage(backgroundImageRaw)
	BoxImage = LoadImage(boxImageRaw)
	HeroImage = LoadImage(hero)
	BulletImage = LoadImage(bullet)
	PlatformImage = LoadImage(platform)
	DefaultWeaponImage = LoadImage(defaultWeapon)
	SniperRifleWeaponImage = LoadImage(sniperRifleWeapon)
	MachineGunWeaponImage = LoadImage(machineGunWeapon)
}
