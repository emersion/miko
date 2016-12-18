// The game engine.
package engine

import (
	"container/list"
	"errors"
	"io"
	"log"
	"time"

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
	"git.emersion.fr/saucisse-royale/miko.git/server/timeserver"
)

const broadcastInterval = 150 * time.Millisecond

type Engine struct {
	auth    *auth.Service
	clock   *clock.Service
	entity  *entity.Service
	action  *action.Service
	terrain *terrain.Terrain
	config  *game.Config

	ctx     *message.Context
	srv     *server.Server
	timeSrv *timeserver.Server

	clients map[int]*server.Client

	mover          *Mover
	running        bool
	brdStop        chan bool
	listenStop     chan bool
	listenTimeStop chan bool
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
		if err := e.entity.AcceptRequest(r); err != nil {
			// This shouldn't fail
			panic(err)
		}

		session := e.auth.GetSessionByEntity(req.Action.Initiator)

		// TODO: error handling
		e.clients[session.Id].Write(func (w io.Writer) error {
			return builder.SendEntityIdChange(w, req.Action.Params[1].(message.EntityId), ball.Id)
		})
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

		if err := e.action.AcceptRequest(r); err != nil {
			log.Println("Warning: Error while accepting action request:", err)
		}
	case entity.Request:
		if !e.processEntityRequest(r) {
			req.Done(errors.New("Request has been rejected"))
			return
		}

		if err := e.entity.AcceptRequest(r); err != nil {
			log.Println("Warning: Error while accepting entity request:", err)
		}
	}
}

func (e *Engine) moveEntities(t message.AbsoluteTick) {
	for _, entity := range e.entity.List() {
		req := e.mover.UpdateEntity(entity, t)
		if req != nil {
			// Let's assume mover already did all security checks
			if err := e.entity.AcceptRequest(req); err != nil {
				// This should not fail
				panic(err)
			}
		}
	}
}

func (e *Engine) broadcastChanges() {
	if e.ctx.Entity.IsDirty() {
		log.Println("Entity diff dirty, broadcasting to clients...")
		err := e.srv.Write(func (w io.Writer) error {
			return builder.SendEntitiesDiffToClients(w, e.clock.GetRelativeTick(), e.ctx.Entity.Flush())
		})
		if err != nil {
			log.Println("Cannot broadcast entities diff:", err)
		}
	}
	if e.ctx.Action.IsDirty() {
		log.Println("Actions diff dirty, broadcasting to clients...")
		err := e.srv.Write(func (w io.Writer) error {
			return builder.SendActionsDone(w, e.clock.GetRelativeTick(), e.ctx.Action.Flush())
		})
		if err != nil {
			log.Println("Cannot broadcast actions:", err)
		}
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
		if t == message.InvalidTick {
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
		log.Println("Back to the past!", currentTick, acceptedMinTick)
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
				err := e.entity.AcceptRequest(entity.NewDeleteRequest(currentTick, ent.Id))
				if err != nil {
					// This shouldn't fail
					panic(err)
				}
			} else {
				ent.Attributes[game.TicksLeftAttr] = ttl - 1
			}
		}
	}

	// Cleanup
	e.entity.Cleanup(currentTick)
	e.action.Cleanup(currentTick)

	// DEBUG: print all entites
	log.Printf("Tick %v finished, entities:\n", currentTick)
	for _, ent := range e.entity.List() {
		log.Printf("entity=%+v position=%+v speed=%+v", ent, ent.Position, ent.Speed)
	}
}

func (e *Engine) listenNewClients() {
	hdlr := handler.New(e.ctx)

	for {
		select {
		case c := <-e.srv.Joins:
			log.Println("Accepting client")
			e.clients[c.Id] = c
			go hdlr.Listen(c.Conn)
		case <-e.listenStop:
			return
		}
	}
}

func (e *Engine) listenNewTimeClients() {
	for {
		select {
		case client := <-e.timeSrv.Joins:
			// TODO: Check that the client is connected
			log.Println("Accepting time server client")
			go e.timeSrv.Accept(client)
		case <-e.listenTimeStop:
			return
		}
	}
}

// Periodically broadcast changes to clients.
func (e *Engine) startBrd() {
	ticker := time.NewTicker(broadcastInterval)
	for {
		select {
		case <-e.brdStop:
			ticker.Stop()
			return
		case <-ticker.C:
			e.broadcastChanges()
		}
	}
}

func (e *Engine) Start() {
	e.running = true

	if e.srv != nil {
		go e.listenNewClients()
		go e.startBrd()
	}
	if e.timeSrv != nil {
		go e.listenNewTimeClients()
	}

	engineStart := time.Duration(time.Now().UnixNano()) * time.Nanosecond

	for {
		e.clock.Tick()
		tick := e.clock.GetAbsoluteTick()

		tickStart := time.Duration(time.Now().UnixNano()) * time.Nanosecond
		e.executeTick(tick)
		tickEnd := time.Duration(time.Now().UnixNano()) * time.Nanosecond

		duration := tickEnd - tickStart
		if clock.TickDuration < duration {
			log.Println("Warning: loop duration exceeds tick duration", duration)
		}

		if !e.running {
			return
		}

		time.Sleep(clock.TickDuration*time.Duration(tick) + engineStart - tickEnd)
	}
}

func (e *Engine) Stop() {
	e.running = false

	if e.srv != nil {
		e.brdStop <- true
		e.listenStop <- true
	}

	for _, client := range e.clients {
		client.Write(func (w io.Writer) error {
			return builder.SendExit(w, message.ExitCodes["server_closed"])
		})
		client.Close()
	}
}

func (e *Engine) Context() *message.Context {
	return e.ctx
}

func New(srv *server.Server, timeSrv *timeserver.Server) *Engine {
	// Create a new context
	ctx := message.NewServerContext()

	// Create the engine
	e := &Engine{
		ctx:     ctx,
		srv:     srv,
		timeSrv: timeSrv,
		auth:    auth.NewService(),
		clock:   clock.NewService(),
		entity:  entity.NewService(),
		action:  action.NewService(),
		terrain: terrain.New(),
		config:  game.DefaultConfig(),

		clients:        make(map[int]*server.Client),
		brdStop:        make(chan bool),
		listenStop:     make(chan bool),
		listenTimeStop: make(chan bool),
	}

	// Set config
	if timeSrv != nil {
		e.config.TimeServerPort = uint16(timeSrv.Port())
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

	// Set callbacks
	e.auth.LoginCallback = func(session *message.Session) {
		entity := session.Entity

		// Populate entity
		// TODO: default values are hardcoded
		entity.Position.BX = 20
		entity.Position.BY = 20
		entity.Type = game.PlayerEntity
		entity.Sprite = game.PlayerSprite
		entity.Attributes[game.HealthAttr] = game.Health(1000)
		entity.Attributes[game.CooldownOneAttr] = game.Cooldown(0)
	}

	return e
}

func NewServerless() *Engine {
	return New(nil, nil)
}
