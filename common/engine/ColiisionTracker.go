package engine

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"github.com/simpletonDL/box2d"
)

type CollisionTracker struct {
	engine *GameEngine
}

func NewCollisionTracker(engine *GameEngine) CollisionTracker {
	return CollisionTracker{engine: engine}
}

func (c CollisionTracker) BeginContact(contact box2d.B2ContactInterface) {
	bodyA := contact.GetFixtureA().GetBody()
	bodyB := contact.GetFixtureB().GetBody()
	userDataA := GetBodyUserData(bodyA)
	userDataB := GetBodyUserData(bodyB)

	if userDataA.GetKind() == protocol.BodyKindHero && userDataB.GetKind() == protocol.BodyKindPlatform ||
		userDataA.GetKind() == protocol.BodyKindPlatform && userDataB.GetKind() == protocol.BodyKindHero {
		// Make sure that bodyA is a hero
		if userDataB.GetKind() == protocol.BodyKindHero {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
		}
		c.processHeroWithPlatformBeginContact(contact, bodyA, userDataA.(PlayerUserData), bodyB)
	}
	if userDataA.GetKind() == protocol.BodyKindHero && bodyB.M_type != box2d.B2BodyType.B2_kinematicBody ||
		bodyA.M_type != box2d.B2BodyType.B2_kinematicBody && userDataB.GetKind() == protocol.BodyKindHero {
		// Make sure that bodyA is a hero
		if userDataB.GetKind() == protocol.BodyKindHero {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
		}
		c.processHeroWithStaticOrDynamicBodyBeginContact(contact, bodyA, userDataA.(PlayerUserData), bodyB)
	}

	if userDataA.GetKind() == protocol.BodyKindHero && userDataB.GetKind() == protocol.BodyKindWeaponBox ||
		userDataA.GetKind() == protocol.BodyKindWeaponBox && userDataB.GetKind() == protocol.BodyKindHero {
		// Make sure that bodyA is a hero
		if userDataB.GetKind() == protocol.BodyKindHero {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
		}
		c.processHeroWithWeaponBoxContact(contact, userDataA.(PlayerUserData), bodyB)
	}
}

func (c CollisionTracker) EndContact(contact box2d.B2ContactInterface) {}

func (c CollisionTracker) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	bodyA := contact.GetFixtureA().GetBody()
	bodyB := contact.GetFixtureB().GetBody()
	userDataA := GetBodyUserData(bodyA)
	userDataB := GetBodyUserData(bodyB)

	if userDataA.GetKind() == protocol.BodyKindBullet || userDataB.GetKind() == protocol.BodyKindBullet {
		// Make sure that bodyA is a bullet
		if userDataB.GetKind() == protocol.BodyKindBullet {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
		}
		c.processBulletPreSolveContact(contact, bodyA, bodyB)
	}
	if userDataA.GetKind() == protocol.BodyKindHero && userDataB.GetKind() == protocol.BodyKindPlatform ||
		userDataA.GetKind() == protocol.BodyKindPlatform && userDataB.GetKind() == protocol.BodyKindHero {
		// Make sure that bodyA is a hero
		if userDataB.GetKind() == protocol.BodyKindHero {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
		}
		c.processHeroWithPlatformPreSolveContact(contact, userDataA.(PlayerUserData))
	}
}

func (c CollisionTracker) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}

func (c CollisionTracker) processBulletPreSolveContact(contact box2d.B2ContactInterface, bulletBody *box2d.B2Body, otherBody *box2d.B2Body) {
	bulletUserData := GetBodyUserData(bulletBody).(BulletUserData)
	otherUserData := GetBodyUserData(otherBody)
	if otherBody == bulletUserData.Owner {
		// bullets should not contact with their owners
		contact.SetEnabled(false)
		return
	}

	fmt.Printf("Bullet(%f, %f) contact\n", bulletBody.GetPosition().X, bulletBody.GetPosition().Y)
	contact.SetEnabled(false)
	c.engine.ScheduleCommand(RemoveBodyCommand{body: bulletBody})
	if otherUserData.GetKind() == protocol.BodyKindBullet {
		// Ignore contact with other bullets
	} else {
		var worldManifold box2d.B2WorldManifold
		contact.GetWorldManifold(&worldManifold)
		if contact.GetManifold().PointCount > 0 {
			collisionPoint := worldManifold.Points[0]
			fmt.Printf("Collision point: (%f, %f)\n", collisionPoint.X, collisionPoint.Y)
			fmt.Printf("Body world center: (%f, %f)", otherBody.GetWorldCenter().X, otherBody.GetWorldCenter().Y)
			bulletVelocity := bulletBody.GetLinearVelocity()
			bulletVelocity.Normalize()
			bulletVelocity.OperatorScalarMulInplace(bulletUserData.ImpactForce)
			c.engine.ScheduleCommand(ApplyImpulseCommand{
				body:    otherBody,
				point:   collisionPoint,
				impulse: bulletVelocity,
			})
		}
	}
}

func (c CollisionTracker) processHeroWithPlatformBeginContact(
	contact box2d.B2ContactInterface, heroBody *box2d.B2Body, heroUserData PlayerUserData, platformBody *box2d.B2Body,
) {
	// Don't allow player to move down through platform in case of contact begin
	c.engine.Players[heroUserData.HeroId].MoveDownThrowPlatform = false

	var woldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&woldManifold)

	platformY := platformBody.GetPosition().Y
	for i := 0; i < contact.GetManifold().PointCount; i++ {
		contactPointY := woldManifold.Points[i].Y
		if contactPointY > platformY {
			// Since this method is called in BeginContact it means that hero first time contact with platform.
			// If contact point is upper that platform center then its mean that we should preserve contact.
			return
		}
	}
	// All contact points are under platform
	contact.SetEnabled(false)
}

func (c CollisionTracker) processHeroWithPlatformPreSolveContact(contact box2d.B2ContactInterface, heroUserData PlayerUserData) {
	playerInfo := c.engine.Players[heroUserData.HeroId]
	if playerInfo.MoveDownThrowPlatform {
		contact.SetEnabled(false)
	}
	playerInfo.MoveDownThrowPlatform = false
}

func (c CollisionTracker) processHeroWithStaticOrDynamicBodyBeginContact(contact box2d.B2ContactInterface, heroBody *box2d.B2Body, heroUserData PlayerUserData, otherBody *box2d.B2Body) {
	otherBodyY := otherBody.GetPosition().Y

	var woldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&woldManifold)

	for i := 0; i < contact.GetManifold().PointCount; i++ {
		contactPointY := woldManifold.Points[i].Y
		if contactPointY < otherBodyY {
			return
		}
	}
	// All contact points are over platform/box/e.t.c.
	playerInfo := c.engine.Players[heroUserData.HeroId]
	playerInfo.JumpCount = settings.PlayerMaxJumpCount
}

func (c CollisionTracker) processHeroWithWeaponBoxContact(contact box2d.B2ContactInterface, heroUserData PlayerUserData, weaponBoxBody *box2d.B2Body) {
	// Set new random weapon
	playerInfo := c.engine.Players[heroUserData.HeroId]
	c.engine.ScheduleCommand(NewCustomCommand(func(engine *GameEngine) {
		kind := utils.RandInRange(1, int(protocol.WeaponKindCount))
		playerInfo.Weapon = CreateWeapon(protocol.WeaponKind(kind))
	}))
	contact.SetEnabled(false)
	c.engine.ScheduleCommand(RemoveBodyCommand{body: weaponBoxBody})
}
