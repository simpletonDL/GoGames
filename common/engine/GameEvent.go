package engine

import (
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/box2d"
	"math/rand"
	"time"
)

type GameEvent interface {
	GetFrequency(engine *GameEngine) time.Duration
	ProcessEvent(engine *GameEngine)
}

type WeaponBoxCreationEvent struct {
	frequency         time.Duration
	boxesCountPerTime int
	currentBoxes      *box2d.B2Body
}

func (w *WeaponBoxCreationEvent) GetFrequency(engine *GameEngine) time.Duration {
	return w.frequency
}

func (w *WeaponBoxCreationEvent) ProcessEvent(engine *GameEngine) {
	x := rand.Uint32() % settings.WorldWidth
	AddWeaponBox(engine.World, float64(x), settings.WorldHeight, rand.Float64(), 1, 1, 1, 0.4)
}

func NewWeaponBoxCreationEvent(frequency time.Duration, boxesCountPerTime int) *WeaponBoxCreationEvent {
	return &WeaponBoxCreationEvent{
		frequency:         frequency,
		boxesCountPerTime: boxesCountPerTime,
		currentBoxes:      nil,
	}
}

func (e *GameEngine) RunEventAsync(event GameEvent) {
	go func() {
		ticker := time.NewTicker(event.GetFrequency(e))
		for {
			<-ticker.C
			e.ScheduleCommand(NewCustomCommand(func(e *GameEngine) {
				event.ProcessEvent(e)
			}))
			ticker.Reset(event.GetFrequency(e))
		}
	}()
}
