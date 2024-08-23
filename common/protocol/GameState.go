package protocol

type BodyKind uint8
type DirectionKind bool
type WeaponKind int

const (
	BodyKindBox = BodyKind(iota)
	BodyKindHero
	BodyKindBullet
	BodyKindPlatform
	BodyKindWeaponBox
)

const (
	DirectionKindRight = DirectionKind(true)
	DirectionKindLeft  = DirectionKind(false)
)

const (
	WeaponKindDefault = WeaponKind(iota)
	WeaponKindSniperRifle
	WeaponKindMachineGun

	WeaponKindCount
)

type GameState struct {
	Objects []GameObject
}

func NewEmptyGameState() GameState {
	return GameState{
		Objects: []GameObject{},
	}
}

type GameObject struct {
	XPos          float64
	YPos          float64
	Angel         float64
	BodyKind      BodyKind
	Width, Height float64
	Direction     DirectionKind
	// Player specific
	WeaponKind                       WeaponKind
	WeaponAvailableBullets           int
	WeaponAvailableBulletsInMagazine int
	WeaponMagazineCapacity           int
	WeaponIsReady                    bool
}
