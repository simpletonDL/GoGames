package settings

import (
	"math"
	"time"
)

const (
	/* Game size parameters */

	WorldWidth   = 24
	WorldHeight  = 16
	ScreenWidth  = 1200
	ScreenHeight = 800

	OutOfScreenBound = 6

	GameInputCapacity = 1000
	GameFPS           = 60

	/* Box2D parameters */

	VelocityIterations = 6
	PositionIterations = 2

	/* Player parameters */

	PlayerMaxHorizontalSpeed        = 8.0
	PlayerGetMaxHorizontalSpeedTime = time.Millisecond * 300

	PlayerJumpSpeed    = 10.0
	PlayerDownSpeed    = -8.0
	PlayerMaxJumpCount = 2
	PlayerAngularSpeed = math.Pi

	PlayerLivesCount = 7
)

/* Inferred parameters */

const (
	PlayerHorizontalAccelerationPerFrame = PlayerMaxHorizontalSpeed /
		(float64(PlayerGetMaxHorizontalSpeedTime) / float64(time.Second) * GameFPS)
)

/* Dev parameters */

const (
	Debug = true
)
