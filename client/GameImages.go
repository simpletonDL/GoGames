package client

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

var (
	//go:embed assets/stars-background.jpg
	backgroundImageRaw []byte

	//go:embed assets/wooden-box.png
	boxImageRaw []byte

	//go:embed assets/weapon-box.png
	weaponBoxImage []byte

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

	//go:embed assets/weapons/сarbine.png
	carbineWeapon []byte
)

var (
	BackgroundImage *ebiten.Image
	BulletImage     *ebiten.Image
	PlatformImage   *ebiten.Image
	BoxImage        *ebiten.Image
	HeroImage       *ebiten.Image
	WeaponBoxImage  *ebiten.Image

	/* Weapons */

	DefaultWeaponImage     *ebiten.Image
	SniperRifleWeaponImage *ebiten.Image
	MachineGunWeaponImage  *ebiten.Image
	CarbineWeaponImage     *ebiten.Image
)

var (
	MainFont *text.GoTextFaceSource
)

func LoadImage(bs []byte) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func loadImages() {
	BackgroundImage = LoadImage(backgroundImageRaw)
	BoxImage = LoadImage(boxImageRaw)
	HeroImage = LoadImage(hero)
	BulletImage = LoadImage(bullet)
	PlatformImage = LoadImage(platform)
	DefaultWeaponImage = LoadImage(defaultWeapon)
	SniperRifleWeaponImage = LoadImage(sniperRifleWeapon)
	MachineGunWeaponImage = LoadImage(machineGunWeapon)
	CarbineWeaponImage = LoadImage(carbineWeapon)
	WeaponBoxImage = LoadImage(weaponBoxImage)
}

func loadFont() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	MainFont = s
}

func LoadImagesAndFonts() {
	loadImages()
	loadFont()
}
