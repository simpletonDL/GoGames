package engine

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/protocol"
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
	userDataA := bodyA.GetUserData().(BodyUserData)
	userDataB := bodyB.GetUserData().(BodyUserData)

	if userDataA.Kind == protocol.BodyKind.Hero && userDataB.Kind == protocol.BodyKind.Platform ||
		userDataA.Kind == protocol.BodyKind.Platform && userDataB.Kind == protocol.BodyKind.Hero {
		// Make sure that bodyA is a hero
		if userDataB.Kind == protocol.BodyKind.Hero {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
			c.processHeroWithPlatformContact(contact, bodyA, bodyB)
		}
	}
}

func (c CollisionTracker) EndContact(contact box2d.B2ContactInterface) {}

func (c CollisionTracker) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	bodyA := contact.GetFixtureA().GetBody()
	bodyB := contact.GetFixtureB().GetBody()
	userDataA := bodyA.GetUserData().(BodyUserData)
	userDataB := bodyB.GetUserData().(BodyUserData)

	if userDataA.Kind == protocol.BodyKind.Bullet || userDataB.Kind == protocol.BodyKind.Bullet {
		// Make sure that bodyA is a bullet
		if userDataB.Kind == protocol.BodyKind.Bullet {
			bodyA, bodyB = bodyB, bodyA
			userDataA, userDataB = userDataB, userDataA
		}
		c.processBulletContact(contact, bodyA, bodyB)
	}
}

func (c CollisionTracker) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}

func (c CollisionTracker) processBulletContact(contact box2d.B2ContactInterface, bulletBody *box2d.B2Body, otherBody *box2d.B2Body) {
	bulletUserData := bulletBody.GetUserData().(BodyUserData)
	otherUserData := otherBody.GetUserData().(BodyUserData)
	if otherBody == bulletUserData.Owner {
		// bullets should not contact with their owners
		contact.SetEnabled(false)
		return
	}

	fmt.Printf("Bullet(%f, %f) contact\n", bulletBody.GetPosition().X, bulletBody.GetPosition().Y)
	contact.SetEnabled(false)
	c.engine.ScheduleCommand(RemoveBodyCommand{body: bulletBody})
	if otherUserData.Kind == protocol.BodyKind.Bullet {
		//c.engine.ScheduleCommand(RemoveBodyCommand{body: otherBody})
	} else {
		var worldManifold box2d.B2WorldManifold
		contact.GetWorldManifold(&worldManifold)
		if contact.GetManifold().PointCount > 0 {
			collisionPoint := worldManifold.Points[0]
			fmt.Printf("Collision point: (%f, %f)\n", collisionPoint.X, collisionPoint.Y)
			fmt.Printf("Body world center: (%f, %f)", otherBody.GetWorldCenter().X, otherBody.GetWorldCenter().Y)
			c.engine.ScheduleCommand(ApplyImpulseCommand{
				body:    otherBody,
				point:   collisionPoint,
				impulse: box2d.B2Vec2{X: 7, Y: 0},
			})
		}
	}
}

func (c CollisionTracker) processHeroWithPlatformContact(contact box2d.B2ContactInterface, heroBody *box2d.B2Body, platformBody *box2d.B2Body) {
	var woldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&woldManifold)

	platformY := platformBody.GetPosition().Y
	for i := 0; i < len(woldManifold.Points); i++ {
		contactPointY := woldManifold.Points[i].Y
		if contactPointY > platformY {
			// Since this method is called in BeginContact its mean that hero first time contact with platform.
			// If contact point is upper that platform center then its mean that we should preserve contact.
			return
		}
	}
	// All contact points are under platform
	contact.SetEnabled(false)
}
