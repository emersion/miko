// The game engine.
package engine

import (
	"container/list"
	"errors"
	"git.emersion.fr/saucisse-royale/miko.git/server/action"
	"git.emersion.fr/saucisse-royale/miko.git/server/auth"
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/game"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
	"log"
	"time"
)

const broadcastInterval time.Duration = time.Millisecond * 150

type Engine struct {
	auth    *auth.Service
	clock   *clock.Service
	entity  *entity.Service
	action  *action.Service
	terrain *terrain.Terrain
	config  *game.Config

	ctx *message.Context
	srv *server.Server

	clients map[int]*message.IO

	mover     *Mover
	ticker    *time.Ticker
	brdTicker *time.Ticker
	stop      chan bool
}

func (e *Engine) processEntityRequest(req entity.Request) bool {
	return true
}

func (e *Engine) processActionRequest(req *action.Request) bool {
	if req.Action.Id == game.ThrowBallAction {
		log.Println("Accepting ball", req)

		initiator := e.entity.Get(req.Action.Initiator)
		ball := entity.New()
		ball.Type = game.BallEntity
		ball.Sprite = game.BallSprite
		ball.Position = initiator.Position
		ball.Speed.Norm = float64(e.config.DefaultBallSpeed)
		ball.Speed.Angle = float64(req.Action.Params[0].(float32))
		ball.Attributes[game.TicksLeftAttr] = e.config.DefaultBallLifespan
		ball.Attributes[game.SenderAttr] = initiator.Id
		r := entity.NewCreateRequest(req.GetTick(), ball)
		e.entity.AcceptRequest(r)

		session := e.auth.GetSessionByEntity(req.Action.Initiator)
		client := e.clients[session.Id]
		builder.SendEntityIdChange(client, req.Action.Params[1].(message.EntityId), ball.Id)
	}

	return true
}

func (e *Engine) processRequest(req message.Request) {
	switch r := req.(type) {
	case *action.Request:
		if !e.processActionRequest(r) {
			req.Done(errors.New("Request has been rejected"))
			return
		}

		err := e.action.AcceptRequest(r)
		if err != nil {
			log.Println("Warning: Error while accepting action request:", err)
		}
	case entity.Request:
		if !e.processEntityRequest(r) {
			req.Done(errors.New("Request has been rejected"))
			return
		}

		err := e.entity.AcceptRequest(r)
		if err != nil {
			log.Println("Warning: Error while accepting entity request:", err)
		}
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
		log.Println("Entity diff dirty, broadcasting to clients...")
		err := builder.SendEntitiesDiffToClients(e.srv, e.clock.GetRelativeTick(), e.ctx.Entity.Flush())
		if err != nil {
			log.Println("Cannot broadcast entities diff:", err)
		}
	}
	if e.ctx.Action.IsDirty() {
		log.Println("Actions diff dirty, broadcasting to clients...")
		err := builder.SendActionsDone(e.srv, e.clock.GetRelativeTick(), e.ctx.Action.Flush())
		if err != nil {
			log.Println("Cannot broadcast actions:", err)
		}
	}
}

func (e *Engine) startBroadcastingChanges() {
	e.brdTicker = time.NewTicker(broadcastInterval)
	for {
		<-e.brdTicker.C
		e.broadcastChanges()
	}
}

func (e *Engine) executeTick(currentTick message.AbsoluteTick) {
	entityFrontend := e.entity.Frontend()
	actionFrontend := e.action.Frontend()

	// Process requests from clients
	minTick := currentTick - message.MaxRewind
	if message.MaxRewind > currentTick {
		minTick = 0
	}
	acceptedMinTick := currentTick
	accepted := list.New()
	for {
		var req message.Request

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
		if t == 0 {
			req.Done(errors.New("Invalid tick"))
			log.Println("Warning: Dropped request, invalid tick")
			continue
		}
		if t < minTick {
			req.Done(errors.New("Request is too old"))
			log.Println("Warning: Dropped request, too old", currentTick-t)
			continue
		}
		if t > currentTick {
			req.Done(errors.New("Request is in the future"))
			log.Println("Warning: Dropped request, in the future", currentTick-t)
			continue
		}
		if t < acceptedMinTick {
			acceptedMinTick = t
		}

		// Append request to the list, keeping it ordered
		inserted := false
		for e := accepted.Front(); e != nil; e = e.Next() {
			r := e.Value.(message.Request)

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

	if accepted.Len() > 0 {
		//log.Println("Accepted", accepted.Len(), "requests from clients")
	}

	// Initiate lag compensation if necessary
	if acceptedMinTick < currentTick {
		//log.Println("Back to the past!", currentTick, currentTick-acceptedMinTick)
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

			// Do not redo deltas not triggered by the user
			// These are the ones that will be computed again
			if !ed.Requested() {
				entDeltas.Remove(entEl)
				entEl = entEl.Next()
				continue
			}

			entEl = entEl.Next()

			// Compute new entities positions
			e.moveEntities(ed.GetTick())

			// Accept new requests at the correct tick
			for el := accepted.Front(); el != nil; el = el.Next() {
				req := el.Value.(message.Request)

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
	e.moveEntities(currentTick)

	// Destroy entities that are not alive anymore
	// TODO: history support
	for _, ent := range e.entity.List() {
		if val, ok := ent.Attributes[game.TicksLeftAttr]; ok {
			ttl := val.(game.TicksLeft)
			if ttl == 0 { // Destroy entity
				e.entity.AcceptRequest(entity.NewDeleteRequest(currentTick, ent.Id))
			} else {
				ent.Attributes[game.TicksLeftAttr] = ttl - 1
			}
		}
	}

	// Cleanup
	e.entity.Cleanup(currentTick)
	e.action.Cleanup(currentTick)
}

func (e *Engine) Start() {
	if e.srv != nil {
		go e.listenNewClients()
		go e.startBroadcastingChanges()
	}

	//e.ticker = time.NewTicker(clock.TickDuration)
	engineStart := time.Duration(time.Now().UnixNano()) * time.Nanosecond

	for {
		// Stop the engine?
		select {
		case <-e.stop:
			return
		default:
		}

		//<-e.ticker.C

		e.clock.Tick()
		tick := e.clock.GetAbsoluteTick()

		tickStart := time.Duration(time.Now().UnixNano()) * time.Nanosecond
		e.executeTick(tick)
		tickEnd := time.Duration(time.Now().UnixNano()) * time.Nanosecond

		duration := tickEnd - tickStart
		if clock.TickDuration < duration {
			log.Println("Warning: loop duration exceeds tick duration", duration)
		}

		log.Println("Time:", (tickStart-engineStart)/time.Duration(tick))
		time.Sleep(clock.TickDuration*time.Duration(tick+1) + engineStart - tickEnd)
	}
}

func (e *Engine) Stop() {
	if e.ticker != nil {
		e.ticker.Stop()
	}
	if e.brdTicker != nil {
		e.brdTicker.Stop()
	}
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
		config:  game.DefaultConfig(),

		clients: make(map[int]*message.IO),
		stop:    make(chan bool),
	}

	// Populate context
	ctx.Auth = e.auth
	ctx.Entity = e.entity.Frontend()
	ctx.Action = e.action.Frontend()
	ctx.Terrain = e.terrain
	ctx.Clock = e.clock
	ctx.Config = e.config

	// Initialize engine submodules
	e.mover = NewMover(e)

	return e
}
