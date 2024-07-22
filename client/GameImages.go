package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

var (
	BackgroundImage *ebiten.Image
	PlayerImage     *ebiten.Image
	BulletImage     *ebiten.Image
	BoxImage        *ebiten.Image
)

func LoadImage(path string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func LoadImages(pathToAssets string) {
	//PlayerImage = LoadImage("assets/submarine.png")
	//BulletImage = LoadImage("assets/bullet.png")
	BackgroundImage = LoadImage("assets/stars-background.jpg")
	BoxImage = LoadImage(pathToAssets + "/wooden-box.png")
}
