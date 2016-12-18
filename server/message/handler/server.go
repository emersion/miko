package handler

import (
	"errors"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"io"
	"log"
)

func ReadVersion(r io.Reader) (version message.ProtocolVersion) {
	Read(r, &version)
	return
}

func ReadLogin(r io.Reader) (username, password string) {
	Read(r, &username, &password)
	return
}

func ReadRegister(r io.Reader) (username, password string) {
	Read(r, &username, &password)
	return
}

func ReadChatSend(r io.Reader) (msg string) {
	Read(r, &msg)
	return
}

var serverHandlers = &map[message.Type]TypeHandler{
	message.Types["version"]: func(ctx *message.Context, conn *message.Conn) error {
		conn.Version = ReadVersion(conn)

		if conn.Version != message.CurrentVersion {
			log.Println("Protocol version mismatch")

			code := message.ExitCodes["client_outdated"]
			if conn.Version > message.CurrentVersion {
				code = message.ExitCodes["server_outdated"]
			}

			err := conn.Write(func (w io.Writer) error {
				return builder.SendExit(w, code)
			})
			if err != nil {
				return err
			}

			if err := conn.Close(); err != nil {
				return err
			}
		} else {
			log.Println("Sending config to client")
			conn.State = message.Accepted

			err := conn.Write(func (w io.Writer) error {
				return builder.SendConfig(w, ctx.Config)
			})
			if err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["exit"]: func(ctx *message.Context, conn *message.Conn) error {
		code := ReadExit(conn)
		log.Println("Client disconnected with code:", code)

		if err := conn.Close(); err != nil {
			return err
		}

		return nil
	},
	message.Types["login"]: func(ctx *message.Context, conn *message.Conn) error {
		username, password := ReadLogin(conn)

		// TODO: do not use time.Now(), use the time at which the current tick has begun
		code := ctx.Auth.Login(conn.Id, username, password)
		err := conn.Write(func (w io.Writer) error {
			return builder.SendLoginResp(w, code, ctx.Clock.GetRelativeTick(), ctx.Clock.GetTickTime())
		})
		if err != nil {
			return err
		}

		log.Println("Client requested login:", username, code)

		if code != message.LoginResponseCodes["ok"] {
			return nil
		}

		conn.State = message.LoggedIn

		// Get entity
		session := ctx.Auth.GetSession(conn.Id)
		if session == nil {
			return errors.New("Cannot get newly logged in user's session")
		}

		// Send initial terrain
		log.Println("Sending initial terrain")
		pos := session.Entity.Position
		radius := message.BlockCoord(8)
		blks := []*message.Block{}
		for i := pos.BX - radius; i <= pos.BX+radius; i++ {
			for j := pos.BY - radius; j <= pos.BY+radius; j++ {
				blk, err := ctx.Terrain.GetBlockAt(i, j)
				if err != nil {
					log.Println("Cannot send initial terrain to client", err)
					continue
				}

				blks = append(blks, blk)
			}
		}
		err = conn.Write(func (w io.Writer) error {
			return builder.SendChunksUpdate(w, ctx.Clock.GetRelativeTick(), blks)
		})
		if err != nil {
			return err
		}

		// Send initial entities
		log.Println("Sending initial entities")
		for _, e := range ctx.Entity.List() {
			if e.Id == session.Entity.Id {
				continue
			}

			err := conn.Write(func (w io.Writer) error {
				return builder.SendEntityCreate(w, ctx.Clock.GetRelativeTick(), e)
			})
			if err != nil {
				return err
			}
		}

		// Send current players
		for _, s := range ctx.Auth.List() {
			if s.Id == session.Id {
				continue
			}

			err := conn.Write(func (w io.Writer) error {
				return builder.SendPlayerJoined(w, ctx.Clock.GetRelativeTick(), s.Entity.Id, s.Username)
			})
			if err != nil {
				return err
			}
		}

		// Register new entity
		req := ctx.Entity.Add(session.Entity, ctx.Clock.GetAbsoluteTick()) // TODO: move this elsewhere?
		err = req.Wait()
		log.Println("Entity registered!")
		if err != nil {
			return err
		}

		// Broadcast new entity to other clients
		log.Println("Flushing entities diff")
		err = conn.Broadcast(func (w io.Writer) error {
			return builder.SendEntitiesDiffToClients(w, ctx.Clock.GetRelativeTick(), ctx.Entity.Flush())
		})
		if err != nil {
			return err
		}

		// Send new entity to this client
		err = conn.Write(func (w io.Writer) error {
			return builder.SendEntityCreate(w, ctx.Clock.GetRelativeTick(), session.Entity)
		})
		if err != nil {
			return err
		}

		// Mark the io as ready
		// It will now receive all broadcasts
		conn.State = message.Ready

		// Broadcast player_joined to everyone (including this client)
		err = conn.Broadcast(func (w io.Writer) error {
			return builder.SendPlayerJoined(w, ctx.Clock.GetRelativeTick(), session.Entity.Id, username)
		})
		if err != nil {
			return err
		}

		return nil
	},
	message.Types["register"]: func(ctx *message.Context, conn *message.Conn) error {
		username, password := ReadRegister(conn)

		code := ctx.Auth.Register(conn.Id, username, password)

		log.Println("Client registered:", username, code)

		return conn.Write(func (w io.Writer) error {
			return builder.SendRegisterResp(w, code)
		})
	},
	message.Types["terrain_request"]: func(ctx *message.Context, conn *message.Conn) error {
		var size uint8
		Read(conn, &size)

		for i := 0; i < int(size); i++ {
			var x, y message.BlockCoord
			Read(conn, &x, &y)

			blk, err := ctx.Terrain.GetBlockAt(x, y)
			if err != nil {
				return err
			}

			err = conn.Write(func (w io.Writer) error {
				return builder.SendChunkUpdate(w, ctx.Clock.GetRelativeTick(), blk)
			})
			if err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["entity_update"]: func(ctx *message.Context, conn *message.Conn) error {
		// TODO: add update initiator as parameter to ctx.Entity.Update()

		t := ctx.Clock.ToAbsoluteTick(readTick(conn))

		entity, diff := ReadEntity(conn)
		log.Printf("Received entity update: tick=%v entity=%+v position=%+v speed=%+v\n", t, entity, entity.Position, entity.Speed)
		ctx.Entity.Update(entity, diff, t)

		return nil
	},
	message.Types["action_do"]: func(ctx *message.Context, conn *message.Conn) error {
		if !ctx.Auth.HasSession(conn.Id) {
			return errors.New("Cannot execute action: remote not logged in")
		}

		session := ctx.Auth.GetSession(conn.Id)

		action := &message.Action{
			Initiator: session.Entity.Id,
		}

		t := ctx.Clock.ToAbsoluteTick(readTick(conn))
		Read(conn, &action.Id)

		// TODO: move action params somewhere else
		// [GAME-SPECIFIC]
		if action.Id == 0 { // throw_ball
			var angle float32
			var tmpId message.EntityId
			Read(conn, &angle, &tmpId)
			log.Println("Received ball action:", angle, tmpId)

			action.Params = []interface{}{angle, tmpId}
		} else {
			log.Println("Client triggered unknown action:", action.Id)
		}

		ctx.Action.Execute(action, t)
		return nil
	},
	message.Types["chat_send"]: func(ctx *message.Context, conn *message.Conn) error {
		msg := ReadChatSend(conn)

		if !ctx.Auth.HasSession(conn.Id) {
			return errors.New("User not authenticated")
		}

		session := ctx.Auth.GetSession(conn.Id)
		username := session.Username

		log.Println("Broadcasting chat message:", username, msg)

		return conn.Broadcast(func (w io.Writer) error {
			return builder.SendChatReceive(w, ctx.Clock.GetRelativeTick(), username, msg)
		})
	},
}
