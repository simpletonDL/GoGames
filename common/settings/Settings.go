package settings

const (
	WorldWidth   = 16
	WorldHeight  = 16
	ScreenWidth  = 800
	ScreenHeight = 800

	ServerInputCapacity = 100
	ServerFPS           = 60
	ServerTimestamp     = float64(1) / ServerFPS
	VelocityIterations  = 6
	PositionIterations  = 2
)
