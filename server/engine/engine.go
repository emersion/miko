// The game engine.
package engine

import (
	"container/list"
	"git.emersion.fr/saucisse-royale/miko.git/server/action"
	"git.emersion.fr/saucisse-royale/miko.git/server/auth"
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko.git/server/requests"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
	"log"
	"time"
)

const broadcastInterval time.Duration = time.Millisecond * 1000

type Engine struct {
	auth    *auth.Service
	clock   *clock.Service
	entity  *entity.Service
	action  *action.Service
	terrain *terrain.Terrain
	ctx     *message.Context
	srv     *server.Server

	clients map[int]*message.IO

	mover *Mover
	stop  chan bool
}

func (e *Engine) processEntityRequest(req entity.Request) bool {
	return true
}

func (e *Engine) processActionRequest(req *action.Request) bool {
	if req.Action.Id == 0 { // throw_ball
		initiator := e.entity.Get(req.Action.Initiator)
		ball := entity.New()
		ball.Type = 1   // ball
		ball.Sprite = 2 // ball
		ball.Position = initiator.Position
		ball.Speed.Norm = float64(e.ctx.Config.DefaultBallSpeed)
		ball.Speed.Angle = float64(req.Action.Params[0].(float32))
		r := entity.NewCreateRequest(req.GetTick(), ball)
		e.entity.AcceptRequest(r)

		session := e.auth.GetSessionByEntity(req.Action.Initiator)
		client := e.clients[session.Id]
		builder.SendEntityIdChange(client, req.Action.Params[1].(message.EntityId), ball.Id)
	}

	return true
}

func (e *Engine) processRequest(req requests.Request) {
	switch r := req.(type) {
	case *action.Request:
		if !e.processActionRequest(r) {
			return
		}
		e.action.AcceptRequest(r)
	case entity.Request:
		if !e.processEntityRequest(r) {
			return
		}
		e.entity.AcceptRequest(r)
	}
}

func (e *Engine) moveEntities(t message.AbsoluteTick) {
	for _, entity := range e.entity.List() {
		req := e.mover.UpdateEntity(entity, t)
		if req != nil {
			// Let's assume mover already did all security checks
			e.entity.AcceptRequest(req)
		}
	}
}

func (e *Engine) listenNewClients() {
	hdlr := handler.New(e.ctx)

	for {
		io := <-e.srv.Joins
		e.clients[io.Id] = io
		go hdlr.Listen(io)
	}
}

func (e *Engine) broadcastChanges() {
	if e.ctx.Entity.IsDirty() {
		err := builder.SendEntitiesDiffToClients(e.srv, e.clock.GetRelativeTick(), e.ctx.Entity.Flush())
		if err != nil {
			log.Println("Cannot broadcast entities diff:", err)
		}
	}
	if e.ctx.Action.IsDirty() {
		err := builder.SendActionsDone(e.srv, e.clock.GetRelativeTick(), e.ctx.Action.Flush())
		if err != nil {
			log.Println("Cannot broadcast actions:", err)
		}
	}
}

func (e *Engine) startBroadcastingChanges() {
	for {
		e.broadcastChanges()
		time.Sleep(broadcastInterval)
	}
}

func (e *Engine) Start() {
	go e.listenNewClients()
	go e.startBroadcastingChanges()

	entityFrontend := e.entity.Frontend()
	actionFrontend := e.action.Frontend()

	for {
		start := time.Now().UnixNano()
		e.clock.Tick()

		// Stop the engine?
		select {
		case <-e.stop:
			break
		default:
		}

		// Process requests from clients
		minTick := e.clock.GetAbsoluteTick() - message.MaxRewind
		if message.MaxRewind > e.clock.GetAbsoluteTick() {
			minTick = 0
		}
		acceptedMinTick := e.clock.GetAbsoluteTick()
		accepted := list.New()
		for {
			var req requests.Request

			// Get last request
			noMore := false
			select {
			case req = <-entityFrontend.Creates:
			case req = <-entityFrontend.Updates:
			case req = <-entityFrontend.Deletes:
			case req = <-actionFrontend.Executes:
			default:
				noMore = true
			}
			if noMore {
				break
			}

			t := req.GetTick()
			if t < minTick {
				continue
			}
			if t < acceptedMinTick {
				acceptedMinTick = t
			}

			// Append request to the list, keeping it ordered
			inserted := false
			for e := accepted.Front(); e != nil; e = e.Next() {
				r := e.Value.(requests.Request)

				if r.GetTick() > req.GetTick() {
					accepted.InsertBefore(req, e)
					inserted = true
					break
				}
			}
			if !inserted {
				accepted.PushBack(req)
			}
		}

		// Initiate lag compensation if necessary
		if acceptedMinTick < e.ctx.Clock.GetAbsoluteTick() {
			log.Println("Back to the past!", e.clock.GetAbsoluteTick(), e.clock.GetAbsoluteTick()-acceptedMinTick)
			e.terrain.Rewind(e.terrain.GetTick() - acceptedMinTick)
			e.entity.Rewind(e.entity.GetTick() - acceptedMinTick)
		}

		// Redo all deltas until now
		entDeltas := e.entity.Deltas()
		trnDeltas := e.terrain.Deltas()
		entEl := entDeltas.FirstAfter(acceptedMinTick)
		trnEl := trnDeltas.FirstAfter(acceptedMinTick)
		var ed *entity.Delta
		var td *terrain.Delta
		for entEl != nil || trnEl != nil {
			if entEl != nil {
				ed = entEl.Value.(*entity.Delta)
			} else {
				ed = nil
			}
			if trnEl != nil {
				td = trnEl.Value.(*terrain.Delta)
			} else {
				td = nil
			}

			// TODO: replace terain history by actions history
			if ed == nil || (td != nil && ed.GetTick() >= td.GetTick()) {
				// Process terrain delta
				trnEl = trnEl.Next()

				// Compute new entities positions just before the terrain change
				e.moveEntities(td.GetTick() - 1)

				// Redo terrain change
				e.terrain.Redo(td)
			}
			if td == nil || (ed != nil && ed.GetTick() <= td.GetTick()) {
				// Process entity delta
				entEl = entEl.Next()

				// Do not redo deltas not triggered by the user
				// These are the ones that will be computed again
				if !ed.Requested() {
					entDeltas.Remove(entEl.Prev())
					continue
				}

				// Compute new entities positions
				e.moveEntities(ed.GetTick())

				// Accept new requests at the correct tick
				for el := accepted.Front(); el != nil; el = el.Next() {
					req := el.Value.(requests.Request)

					if req.GetTick() >= ed.GetTick() {
						break
					}

					e.processRequest(req)
					accepted.Remove(el)
				}

				e.processRequest(ed.Request())
			}
		}

		// Accept requests whose tick is > last delta tick
		for el := accepted.Front(); el != nil; el = el.Next() {
			req := el.Value.(entity.Request)
			e.moveEntities(req.GetTick())
			e.processRequest(req)
		}

		// Compute new entities positions
		e.moveEntities(e.clock.GetAbsoluteTick())

		end := time.Now().UnixNano()
		time.Sleep(clock.TickDuration - time.Nanosecond*time.Duration(end-start))
	}
}

func (e *Engine) Stop() {
	e.stop <- true
}

func (e *Engine) Context() *message.Context {
	return e.ctx
}

func New(srv *server.Server) *Engine {
	// Create a new context
	ctx := message.NewServerContext()

	// Create the engine
	e := &Engine{
		ctx:     ctx,
		srv:     srv,
		auth:    auth.NewService(),
		clock:   clock.NewService(),
		entity:  entity.NewService(),
		action:  action.NewService(),
		terrain: terrain.New(),

		clients: make(map[int]*message.IO),
		stop:    make(chan bool),
	}

	// Populate context
	ctx.Auth = e.auth
	ctx.Entity = e.entity.Frontend()
	ctx.Action = e.action.Frontend()
	ctx.Terrain = e.terrain
	ctx.Clock = e.clock

	ctx.Config = &message.Config{
		MaxRollbackTicks:    uint16(message.MaxRewind),
		DefaultPlayerSpeed:  7,
		PlayerBallCooldown:  20,
		DefaultBallSpeed:    9,
		DefaultBallLifespan: 100,
	}

	// Initialize engine submodules
	e.mover = NewMover(e)

	return e
}
