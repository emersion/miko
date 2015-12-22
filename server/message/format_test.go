package message_test

import (
	"bytes"
	"io"
	"testing"

	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
)

type mockCloser struct{}

func (c *mockCloser) Close() (err error) {
	return
}

func newMockIO() *message.IO {
	b := &bytes.Buffer{}
	c := &mockCloser{}
	return message.NewIO(0, b, b, c, nil)
}

type sender func(*message.IO) error

type receiver func(*message.IO, *testing.T)

func testMessage(t *testing.T, msgType message.Type, send sender, receive receiver) {
	io := newMockIO()

	err := send(io)
	if err != nil {
		t.Fatal("Cannot send message:", err)
	}

	receivedType := handler.ReadType(io)
	if msgType != receivedType {
		t.Fatal("Sent type", msgType, "but received", receivedType)
	}

	receive(io, t)

	// Check that there is no more data waiting to be consumed
	remaining := make([]byte, 1)
	n, _ := io.Read(remaining)
	if n != 0 {
		t.Fatal("Some data hasn't been consumed by receiver")
	}
}

func TestPing(t *testing.T) {
	testMessage(t, message.Types["ping"], func(io *message.IO) error {
		return builder.SendPing(io)
	}, func(io *message.IO, t *testing.T) {})
}

func TestPong(t *testing.T) {
	testMessage(t, message.Types["pong"], func(io *message.IO) error {
		return builder.SendPong(io)
	}, func(io *message.IO, t *testing.T) {})
}

func TestExit(t *testing.T) {
	code := message.ExitCodes["client_banned"]

	testMessage(t, message.Types["exit"], func(io *message.IO) error {
		return builder.SendExit(io, code)
	}, func(io *message.IO, t *testing.T) {
		c := handler.ReadExit(io)
		if c != code {
			t.Fatal("Sent code", code, "but received", c)
		}
	})
}

func TestLoginResponse_Ok(t *testing.T) {
	code := message.LoginResponseCodes["ok"]
	tick := message.Tick(42)

	testMessage(t, message.Types["login_response"], func(io *message.IO) error {
		return builder.SendLoginResp(io, code, tick)
	}, func(io *message.IO, t *testing.T) {
		receivedTick, receivedCode := handler.ReadLoginResponse(io)
		if receivedCode != code {
			t.Fatal("Sent code", code, "but received", receivedCode)
		}
		if receivedTick != tick {
			t.Fatal("Sent tick", tick, "but received", receivedTick)
		}
	})
}

func TestLoginResponse_WrongPassword(t *testing.T) {
	code := message.LoginResponseCodes["wrong_password"]

	testMessage(t, message.Types["login_response"], func(io *message.IO) error {
		return builder.SendLoginResp(io, code, 0)
	}, func(io *message.IO, t *testing.T) {
		_, c := handler.ReadLoginResponse(io)
		if c != code {
			t.Fatal("Sent code", code, "but received", c)
		}
	})
}

func TestRegisterResponse(t *testing.T) {
	code := message.RegisterResponseCodes["too_many_tries"]

	testMessage(t, message.Types["register_response"], func(io *message.IO) error {
		return builder.SendRegisterResp(io, code)
	}, func(io *message.IO, t *testing.T) {
		c := handler.ReadRegisterResponse(io)
		if c != code {
			t.Fatal("Sent code", code, "but received", c)
		}
	})
}

func TestMetaAction_PlayerJoined(t *testing.T) {
	code := message.MetaActionCodes["player_joined"]
	tick := message.Tick(42)
	entityId := message.EntityId(69)
	username := "délthàs"

	testMessage(t, message.Types["meta_action"], func(io *message.IO) error {
		return builder.SendPlayerJoined(io, tick, entityId, username)
	}, func(io *message.IO, t *testing.T) {
		rt, id, c, u := handler.ReadMetaAction(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}
		if id != entityId {
			t.Fatal("Sent entity id", entityId, "but received", id)
		}
		if c != code {
			t.Fatal("Sent code", code, "but received", c)
		}
		if u != username {
			t.Fatal("Sent username", username, "but received", u)
		}
	})
}

func TestMetaAction_PlayerLeft(t *testing.T) {
	code := message.MetaActionCodes["player_left"]
	tick := message.Tick(42)
	entityId := message.EntityId(69)

	testMessage(t, message.Types["meta_action"], func(io *message.IO) error {
		return builder.SendPlayerLeft(io, tick, entityId)
	}, func(io *message.IO, t *testing.T) {
		rt, id, c, _ := handler.ReadMetaAction(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}
		if id != entityId {
			t.Fatal("Sent entity id", entityId, "but received", id)
		}
		if c != code {
			t.Fatal("Sent code", code, "but received", c)
		}
	})
}

func TestChunkUpdate(t *testing.T) {
	tick := message.Tick(42)
	blk := message.NewBlock()
	for i := 0; i < 10; i++ {
		blk.Points[10][5+i] = 1
	}
	for i := 0; i < 7; i++ {
		blk.Points[20+i][10] = 2
	}

	testMessage(t, message.Types["chunk_update"], func(io *message.IO) error {
		return builder.SendChunkUpdate(io, tick, blk)
	}, func(io *message.IO, t *testing.T) {
		rt, rblk := handler.ReadChunkUpdate(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}

		for i := range blk.Points {
			for j := range blk.Points[i] {
				pt := blk.Points[i][j]
				rpt := rblk.Points[i][j]
				if rpt != pt {
					t.Fatal("Sent point", pt, " at ", i, j, "but received", rpt)
				}
			}
		}
	})
}

func TestChunksUpdate(t *testing.T) {
	tick := message.Tick(42)

	blk1 := message.NewBlock()
	for i := 0; i < 10; i++ {
		blk1.Points[10][5+i] = 1
	}
	for i := 0; i < 7; i++ {
		blk1.Points[20+i][10] = 2
	}

	blk2 := message.NewBlock()
	blk2.Fill(1)
	blk2.Points[10][5] = 0
	blk2.Points[7][21] = 2

	blks := []*message.Block{blk1, blk2}

	testMessage(t, message.Types["chunks_update"], func(io *message.IO) error {
		return builder.SendChunksUpdate(io, tick, blks)
	}, func(io *message.IO, t *testing.T) {
		rt, rblks := handler.ReadChunksUpdate(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}

		if len(rblks) != len(blks) {
			t.Fatal("Sent", len(blks), "blocks, but received", len(rblks))
		}
		for k, blk := range blks {
			rblk := rblks[k]
			for i := range blk.Points {
				for j := range blk.Points[i] {
					pt := blk.Points[i][j]
					rpt := rblk.Points[i][j]
					if rpt != pt {
						t.Fatal("Sent point", pt, " at ", i, j, "in block", k, "but received", rpt)
					}
				}
			}
		}
	})
}

func TestEntityCreate(t *testing.T) {
	tick := message.Tick(42)
	entity := message.NewEntity()
	entity.Id = 69
	entity.Position.X = 22
	entity.Position.Y = 67
	entity.Speed.Angle = 67
	entity.Speed.Norm = 78
	entity.Sprite = 12
	// TODO: populate other attributes

	testMessage(t, message.Types["entity_create"], func(io *message.IO) error {
		return builder.SendEntityCreate(io, tick, entity)
	}, func(io *message.IO, t *testing.T) {
		rt, e := handler.ReadEntityCreate(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}
		if !entity.Equals(e) {
			t.Fatal("Sent entity", entity, "but received", e)
		}
	})
}

func TestEntitiesUpdate(t *testing.T) {
	tick := message.Tick(42)

	entity1 := message.NewEntity()
	entity1.Id = 69
	entity1.Position.X = 22
	entity1.Position.Y = 65
	diff1 := message.NewEntityDiff()
	diff1.Position = true

	entity2 := message.NewEntity()
	entity2.Id = 76
	entity2.Speed.Angle = 67
	entity2.Speed.Norm = 43
	diff2 := message.NewEntityDiff()
	diff2.SpeedAngle = true
	diff2.SpeedNorm = true

	entities := []*message.Entity{entity1, entity2}
	diffs := []*message.EntityDiff{diff1, diff2}

	testMessage(t, message.Types["entities_update"], func(io *message.IO) error {
		return builder.SendEntitiesUpdate(io, tick, entities, diffs)
	}, func(io *message.IO, t *testing.T) {
		rt, rentities, rdiffs := handler.ReadEntitiesUpdate(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}

		if len(entities) != len(rentities) {
			t.Fatal("Sent", len(entities), "entities but received", len(rentities))
		}
		for i, entity := range entities {
			diff := diffs[i]
			rentity := rentities[i]
			rdiff := rdiffs[i]

			if !diff.Equals(rdiff) {
				t.Fatal("Sent diff", diff, " at offset", i, "but received", rdiff)
			}
			if !entity.EqualsWithDiff(rentity, diff) {
				t.Fatal("Sent entity", entity, " at offset", i, "with diff", diff, "but received", rentity)
			}
		}
	})
}

func TestEntityDestroy(t *testing.T) {
	tick := message.Tick(42)
	entityId := message.EntityId(93)

	testMessage(t, message.Types["entity_destroy"], func(io *message.IO) error {
		return builder.SendEntityDestroy(io, tick, entityId)
	}, func(io *message.IO, t *testing.T) {
		rt, id := handler.ReadEntityDestroy(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}
		if entityId != id {
			t.Fatal("Sent entity id", entityId, "but received", id)
		}
	})
}

func TestActionsDone(t *testing.T) {
	tick := message.Tick(42)

	action1 := message.NewAction()
	action1.Id = 25
	action1.Initiator = 65

	actions := []*message.Action{action1}
	// TODO: more actions tests, with action params

	testMessage(t, message.Types["actions_done"], func(io *message.IO) error {
		return builder.SendActionsDone(io, tick, actions)
	}, func(io *message.IO, t *testing.T) {
		rt, ractions := handler.ReadActionsDone(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}

		if len(actions) != len(ractions) {
			t.Fatal("Sent", len(actions), "actions but received", len(ractions))
		}
		for i, action := range actions {
			raction := ractions[i]
			if !action.Equals(raction) {
				t.Fatal("Sent action", action, "but received", raction)
			}
		}
	})
}

func TestChatReceive(t *testing.T) {
	tick := message.Tick(42)
	username := "délthàs-s@@s"
	msg := "How r u guys? aéoaéoaèo"

	testMessage(t, message.Types["chat_receive"], func(io *message.IO) error {
		return builder.SendChatReceive(io, tick, username, msg)
	}, func(io *message.IO, t *testing.T) {
		rt, rusername, rmsg := handler.ReadChatReceive(io)
		if rt != tick {
			t.Fatal("Sent tick", tick, "but received", rt)
		}
		if rusername != username {
			t.Fatal("Sent username", username, "but received", rusername)
		}
		if rmsg != msg {
			t.Fatal("Sent message", msg, "but received", rmsg)
		}
	})
}

type config struct {
	Int    int64
	Float  float64
	String string
}

func (c *config) WriteTo(w io.Writer) (n int64, err error) {
	err = builder.Write(w, c.Int, c.Float, c.String)
	return
}

func (c *config) ReadFrom(r io.Reader) (n int64, err error) {
	err = handler.Read(r, &c.Int, &c.Float, &c.String)
	return
}

func TestConfig(t *testing.T) {
	conf := &config{42, 69.97543, "Hello World!"}

	testMessage(t, message.Types["config"], func(io *message.IO) error {
		return builder.SendConfig(io, conf)
	}, func(io *message.IO, t *testing.T) {
		rconf := &config{}
		err := handler.ReadConfig(io, rconf)
		if err != nil {
			t.Fatal("Error while reading config:", err)
		}

		if conf.Int != rconf.Int {
			t.Fatal("Sent config item Int with value", conf.Int, "but received", rconf.Int)
		}
		if conf.Float != rconf.Float {
			t.Fatal("Sent config item Float with value", conf.Float, "but received", rconf.Float)
		}
		if conf.String != rconf.String {
			t.Fatal("Sent config item String with value", conf.String, "but received", rconf.String)
		}
	})
}

func TestEntityIdChange(t *testing.T) {
	oldId := message.EntityId(24)
	newId := message.EntityId(87)

	testMessage(t, message.Types["entity_id_change"], func(io *message.IO) error {
		return builder.SendEntityIdChange(io, oldId, newId)
	}, func(io *message.IO, t *testing.T) {
		roldId, rnewId := handler.ReadEntityIdChange(io)
		if roldId != oldId {
			t.Fatal("Sent old id", oldId, "but received", roldId)
		}
		if rnewId != newId {
			t.Fatal("Sent new id", newId, "but received", rnewId)
		}
	})
}
