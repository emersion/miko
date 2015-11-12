package client

import (
	"io"
	"log"
	"math"
	"time"

	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"github.com/gopherjs/gopherjs/js"
)

type EngineInput struct {
	Dirty bool
	Keys  map[string]bool
}

func (i *EngineInput) SetKey(key string, value bool) {
	if i.Keys[key] != value {
		i.Keys[key] = value
		i.Dirty = true
	}
}

func (i *EngineInput) IsKeyActive(key string) bool {
	if val, ok := i.Keys[key]; ok {
		return val
	}
	return false
}

func (i *EngineInput) HandleKeyboardEvent(event *js.Object) {
	value := true
	if event.Get("type").String() == "keyup" {
		value = false
	}

	key := event.Get("key").String()
	if key == "undefined" {
		// See https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/keyCode
		switch event.Get("keyCode").Int() {
		case 40:
			key = "ArrowDown"
		case 37:
			key = "ArrowLeft"
		case 39:
			key = "ArrowRight"
		case 38:
			key = "ArrowUp"
		}
	}

	i.SetKey(key, value)
}

func (i *EngineInput) GetSpeedNorm() float32 {
	keys := []string{"ArrowUp", "ArrowRight", "ArrowLeft", "ArrowDown"}
	for _, key := range keys {
		if i.IsKeyActive(key) {
			return 10
		}
	}
	return 0
}

func (i *EngineInput) GetSpeedAngle() float32 {
	switch {
	case i.IsKeyActive("ArrowRight"):
		return 0
	case i.IsKeyActive("ArrowDown"):
		return math.Pi / 2
	case i.IsKeyActive("ArrowLeft"):
		return math.Pi
	case i.IsKeyActive("ArrowUp"):
		return 3 * math.Pi / 2
	}
	return 0
}

func NewEngineInput() *EngineInput {
	return &EngineInput{Keys: map[string]bool{}}
}

type Engine struct {
	Input *EngineInput
	ctx   *message.Context
	w     io.Writer
}

func (e *Engine) Start() {
	//mover := entity.NewMover(e.ctx.Terrain, e.ctx.Clock)

	var step func(timestampObj *js.Object)
	var lastTick time.Duration
	step = func(timestampObj *js.Object) {
		now := time.Duration(timestampObj.Float()*1000) * time.Microsecond // Precision is 1 µs
		if lastTick == time.Duration(0) {
			lastTick = now
		}

		missedTicks := int((now - lastTick) / clock.TickDuration)
		for i := 0; i < missedTicks; i++ {
			e.ctx.Clock.Tick()
			lastTick += clock.TickDuration
		}

		if e.ctx.Me.Entity != nil && e.Input.Dirty {
			e.Input.Dirty = false

			speed := &message.Speed{
				Norm:  e.Input.GetSpeedNorm(),
				Angle: e.Input.GetSpeedAngle(),
			}
			diff := &message.EntityDiff{SpeedAngle: true, SpeedNorm: true}

			if speed.Norm == 0 {
				diff.Position = true
			}

			e.ctx.Me.Entity.Speed = speed
			e.ctx.Entity.Update(e.ctx.Me.Entity, diff, e.ctx.Clock.GetAbsoluteTick())
		}

		if e.ctx.Entity.IsDirty() {
			err := builder.SendEntitiesDiffToServer(e.w, e.ctx.Clock.GetRelativeTick(), e.ctx.Entity.Flush())
			if err != nil {
				log.Println("Could not send entities update to server", err)
			}
		}

		/*for _, entity := range e.ctx.Entity.List() {
			diff := mover.UpdateEntity(entity)
			if diff != nil {
				e.ctx.Entity.Update(entity, diff, e.ctx.Clock.GetAbsoluteTick())
			}
		}
		e.ctx.Entity.Flush()*/

		js.Global.Call("requestAnimationFrame", step)
	}

	js.Global.Get("document").Call("addEventListener", "keydown", func(event *js.Object) {
		event.Call("preventDefault")
		e.Input.HandleKeyboardEvent(event)
	})
	js.Global.Get("document").Call("addEventListener", "keyup", func(event *js.Object) {
		event.Call("preventDefault")
		e.Input.HandleKeyboardEvent(event)
	})

	js.Global.Call("requestAnimationFrame", step)
}

func NewEngine(ctx *message.Context, w io.Writer) *Engine {
	return &Engine{NewEngineInput(), ctx, w}
}
