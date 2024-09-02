package engine

import (
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"math/rand"
	"time"
)

type GameEvent interface {
	GetFrequency(engine *GameEngine) time.Duration
	ProcessEvent(engine *GameEngine)
}

type ObjectCreationEvent struct {
	frequency          time.Duration
	objectCountPerTime int
	createObjectFun    func(engine *GameEngine)
}

type WeaponBoxCreationEvent struct {
	ObjectCreationEvent
}

func (w *ObjectCreationEvent) GetFrequency(engine *GameEngine) time.Duration {
	return w.frequency
}

func (w *ObjectCreationEvent) ProcessEvent(engine *GameEngine) {
	w.createObjectFun(engine)
}

func NewWeaponBoxCreationEvent(frequency time.Duration, boxesCountPerTime int) *ObjectCreationEvent {
	return &ObjectCreationEvent{
		frequency:          frequency,
		objectCountPerTime: boxesCountPerTime,
		createObjectFun: func(engine *GameEngine) {
			for i := 0; i < boxesCountPerTime; i++ {
				x := rand.Uint32() % settings.WorldWidth
				AddWeaponBox(engine.World, float64(x), settings.WorldHeight, rand.Float64(), 1, 1, 1, 0.4)
			}
		},
	}
}

func NewBoxCreationEvent(frequency time.Duration, boxesCountPerTime int) *ObjectCreationEvent {
	return &ObjectCreationEvent{
		frequency:          frequency,
		objectCountPerTime: boxesCountPerTime,
		createObjectFun: func(engine *GameEngine) {
			for i := 0; i < boxesCountPerTime; i++ {
				x := rand.Uint32() % settings.WorldWidth
				AddBox(engine.World, float64(x), settings.WorldHeight, rand.Float64(), 1, 1, 0.5, 0.3)
			}
		},
	}
}

func (e *GameEngine) RunEvent(event GameEvent) {
	ticker := time.NewTicker(event.GetFrequency(e))
	for {
		<-ticker.C
		select {
		case <-e.Ctx.Done():
			utils.Log("GameEngine Box2D Event cancelled\n")
			return
		default:
		}
		e.ScheduleCommand(NewCustomCommand(func(e *GameEngine) {
			event.ProcessEvent(e)
		}))
		ticker.Reset(event.GetFrequency(e))
	}
}
