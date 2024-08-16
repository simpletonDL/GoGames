package settings

import (
	"math"
	"time"
)

const (
	/* Game size parameters */

	WorldWidth   = 16
	WorldHeight  = 16
	ScreenWidth  = 800
	ScreenHeight = 800

	GameInputCapacity = 1000
	GameFPS           = 60

	/* Box2D parameters */

	VelocityIterations = 6
	PositionIterations = 2

	/* Player parameters */

	PlayerMaxHorizontalSpeed        = 8.0
	PlayerGetMaxHorizontalSpeedTime = time.Millisecond * 300

	PlayerJumpSpeed = 10.0
	PlayerDownSpeed = -7.0

	PlayerAngularSpeed = math.Pi
)

/* Inferred parameters */

const (
	PlayerHorizontalAccelerationPerFrame = PlayerMaxHorizontalSpeed /
		(float64(PlayerGetMaxHorizontalSpeedTime) / float64(time.Second) * GameFPS)
)
