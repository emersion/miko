package client

import (
	"github.com/gopherjs/gopherjs/js"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko/server/browser/client"
)

func NewTerrain(el *js.Object) *client.Terrain {
	t := client.NewTerrain(el)
	t.Reset(1)

	el.Call("addEventListener", "click", func(event *js.Object) {
		x := int(event.Get("clientX").Int() / 5)
		y := int(event.Get("clientY").Int() / 5)

		if t.Points[x][y] == message.PointType(0) {
			t.Points[x][y] = message.PointType(1)
		} else {
			t.Points[x][y] = message.PointType(0)
		}

		t.DrawPoint(x, y)
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
