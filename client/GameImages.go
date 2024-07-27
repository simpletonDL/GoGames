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

	//go:embed assets/artem_aleksyuk.png
	hero []byte
)

var (
	BackgroundImage *ebiten.Image
	PlayerImage     *ebiten.Image
	BulletImage     *ebiten.Image
	BoxImage        *ebiten.Image
	HeroImage       *ebiten.Image
)

func LoadImage(bs []byte) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func LoadImages(pathToAssets string) {
	BackgroundImage = LoadImage(backgroundImageRaw)
	BoxImage = LoadImage(boxImageRaw)
	HeroImage = LoadImage(hero)
}
