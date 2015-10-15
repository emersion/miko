package client

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/emersion/go-js-canvas"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
)

const res = 5

type Terrain struct {
	terrain.Terrain
	canvas *canvas.Canvas
}

func (t *Terrain) drawPoint(x, y int) {
	ptType := t.Points[x][y]

	if int(ptType) == 0 {
		t.canvas.ClearRect(x * res, y * res, res, res)
	} else {
		t.canvas.FillRect(x * res, y * res, res, res)
	}
}

func (t *Terrain) Draw() {
	for i := range t.Points {
		for j := range t.Points[i] {
			t.drawPoint(i, j)
		}
	}
}

func (t *Terrain) SetBlock(blk *message.Block) {
	t.Terrain.SetBlock(blk)
	t.Draw()
}

func NewTerrain() *Terrain {
	el := js.Global.Get("document").Call("getElementById", "canvas")
	el.Set("width", message.BLOCK_LEN * res)
	el.Set("height", message.BLOCK_LEN * res)

	t := &Terrain{}
	t.canvas = canvas.New(el)
	t.Reset(1)

	el.Call("addEventListener", "click", func(event *js.Object) {
		x := int(event.Get("clientX").Int() / res)
		y := int(event.Get("clientY").Int() / res)

		if t.Points[x][y] == message.PointType(0) {
			t.Points[x][y] = message.PointType(1)
		} else {
			t.Points[x][y] = message.PointType(0)
		}

		t.drawPoint(x, y)
	})

	saveBtn := js.Global.Get("document").Call("getElementById", "save-btn")
	saveBtn.Call("addEventListener", "click", func() {
		w := &ExportWriter{}
		builder.WriteBlock(w, t.GetBlockAt(0, 0))
		w.Export()
	})

	loadInput := js.Global.Get("document").Call("getElementById", "load-input")
	loadInput.Call("addEventListener", "change", func() {
		file := loadInput.Get("files").Get("0")
		if file == nil {
			return
		}

		go (func () {
			r := &ExportReader{}
			r.File = file
			blk := handler.ReadBlock(r)
			t.SetBlock(blk)
			t.Draw()
		})()
	})

	return t
}
