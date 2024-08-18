package engine

import (
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/box2d"
	"time"
)

type Weapon interface {
	Shoot(e *GameEngine, playerInfo *PlayerInfo)
	ProcessGameTick()
	GetKind() protocol.WeaponKind
}

type DefaultWeapon struct {
	kind             protocol.WeaponKind
	availableBullets int64

	/* maximum bullet count before reload */
	magazineCapacity           int64
	availableBulletsInMagazine int64

	bulletForce float64
	bulletSpeed float64
	recoilSpeed float64

	/* Reload time (in fps) when magazine is empty */
	reloadTimeFps int
	/* Minimum time (in fps) between two shoots */
	betweenTwoShootsTimeFps int

	/* Time (in fps) is needed to shoot be available. Should be more than remainingTimeToReload */
	remainingTimeToReload   int
	remainingTimeToShootFps int
}

func (c *DefaultWeapon) GetKind() protocol.WeaponKind {
	return c.kind
}

func (c *DefaultWeapon) decrementBullets() {
	c.availableBullets -= 1
	c.availableBulletsInMagazine -= 1
}

func (c *DefaultWeapon) isAvailable() bool {
	return !(c.availableBullets == 0 || c.availableBulletsInMagazine == 0 || c.remainingTimeToShootFps != 0 || c.remainingTimeToReload != 0)
}

func (c *DefaultWeapon) Shoot(engine *GameEngine, playerInfo *PlayerInfo) {
	if !c.isAvailable() {
		return
	}
	c.decrementBullets()
	if c.availableBulletsInMagazine == 0 {
		println("Reload!!!")
		println(c.reloadTimeFps)
		// initiate reload
		c.remainingTimeToReload = c.reloadTimeFps
		println(c.remainingTimeToReload)
		c.remainingTimeToShootFps = 0
	} else {
		// initiate time between shoots
		c.remainingTimeToShootFps = c.betweenTwoShootsTimeFps
		c.remainingTimeToReload = 0
	}

	playerBody := playerInfo.Body
	playerPosition := playerBody.GetPosition()
	// TODO: bullet configuration
	bullet := AddBullet(engine.World, playerPosition.X, playerPosition.Y, 0, 0.2, 0.2, playerBody, c.bulletForce)
	bulletRotation := box2d.MakeB2RotFromAngle(playerBody.GetAngle())
	bulletVec := box2d.MakeB2Vec2(bulletRotation.C, bulletRotation.S)
	bulletVec.OperatorScalarMulInplace(c.bulletSpeed)
	if playerInfo.Direction == protocol.DirectionKindLeft {
		bulletVec.OperatorScalarMulInplace(-1.0)
	}
	bullet.SetLinearVelocity(bulletVec)

	// Process recoilImpulse
	recoilImpulse := bulletVec.Clone()
	recoilImpulse.Normalize()
	recoilImpulse.OperatorScalarMulInplace(-c.recoilSpeed * playerBody.GetMass())
	engine.ScheduleCommand(ApplyImpulseCommand{
		body:    playerBody,
		point:   playerPosition,
		impulse: recoilImpulse,
	})
}

func (c *DefaultWeapon) ProcessGameTick() {
	isReloadFinished := c.remainingTimeToReload == 1
	c.remainingTimeToReload = max(0, c.remainingTimeToReload-1)
	c.remainingTimeToShootFps = max(0, c.remainingTimeToShootFps-1)
	if isReloadFinished {
		c.availableBulletsInMagazine = min(c.magazineCapacity, c.availableBullets)
	}
}

func NewDefaultWeapon(kind protocol.WeaponKind, availableBullets int64, magazineCapacity int64, bulletForce float64, bulletSpeed float64, recoilSpeed float64, reloadTime time.Duration, betweenTwoShootsTime time.Duration) *DefaultWeapon {

	reloadTimeFps := int(float64(reloadTime) / float64(time.Second) * settings.GameFPS)
	println(reloadTimeFps)
	betweenTwoShootsTimeFps := int(float64(betweenTwoShootsTime) / float64(time.Second) * settings.GameFPS)
	return &DefaultWeapon{
		kind:                       kind,
		availableBullets:           availableBullets,
		magazineCapacity:           magazineCapacity,
		availableBulletsInMagazine: min(availableBullets, magazineCapacity),
		bulletForce:                bulletForce,
		bulletSpeed:                bulletSpeed,
		recoilSpeed:                recoilSpeed,
		reloadTimeFps:              reloadTimeFps,
		betweenTwoShootsTimeFps:    betweenTwoShootsTimeFps,
		remainingTimeToReload:      0,
		remainingTimeToShootFps:    0,
	}
}

/* Available weapons */

func NewDefaultGun() Weapon {
	return NewDefaultWeapon(protocol.WeaponKindDefault, 9223372036854775807, 10, 10, 15, 3, time.Second, 200*time.Millisecond)
}

func NewSniperRifle() Weapon {
	return NewDefaultWeapon(protocol.WeaponKindSniperRifle, 20, 1, 30, 30, 20, time.Second, 0)
}
