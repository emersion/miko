package client

import (
	"github.com/gopherjs/gopherjs/js"

	"git.emersion.fr/saucisse-royale/miko.git/server/browser/client"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
)

const res = 5

func getMouseCoords(event *js.Object) (int, int) {
	x := int(event.Get("pageX").Int() / res)
	y := int(event.Get("pageY").Int() / res)
	return x, y
}

func invertPoint(pt message.PointType) message.PointType {
	if pt == message.PointType(0) {
		return message.PointType(1)
	} else {
		return message.PointType(0)
	}
}

func NewTerrain(el *js.Object) *client.Terrain {
	size := 1

	t := client.NewTerrain(el)
	t.Reset(size)
	el.Set("width", size*message.BlockLen*res)
	el.Set("height", size*message.BlockLen*res)

	var pressing bool
	var fromX, fromY, lastX, lastY int
	el.Call("addEventListener", "mousedown", func(event *js.Object) {
		fromX, fromY = getMouseCoords(event)
		lastX, lastY = fromX, fromY
		pressing = true
	})
	el.Call("addEventListener", "mousemove", func(event *js.Object) {
		if !pressing {
			return
		}

		x, y := getMouseCoords(event)
		t.DrawRegion(fromX, fromY, lastX-fromX, lastY-fromY)
		t.Canvas.SetFillStyle("rgba(0,0,255,0.1)")
		t.Canvas.FillRect(fromX*res, fromY*res, (x-fromX)*res, (y-fromY)*res)
		lastX, lastY = x, y
	})
	el.Call("addEventListener", "mouseup", func(event *js.Object) {
		t.DrawRegion(fromX, fromY, lastX-fromX, lastY-fromY)
		pressing = false

		if fromX == lastX && fromY == lastY {
			x, y := fromX, fromY
			pt, _ := t.GetPointAt(x, y)
			t.SetPointAt(x, y, invertPoint(pt), 0)
			t.DrawPoint(x, y)
		} else {
			if fromX > lastX {
				fromX, lastX = lastX, fromX
			}
			if fromY > lastY {
				fromY, lastY = lastY, fromY
			}
			for i := fromX; i < lastX; i++ {
				for j := fromY; j < lastY; j++ {
					pt, _ := t.GetPointAt(i, j)
					t.SetPointAt(i, j, invertPoint(pt), 0)
					t.DrawPoint(i, j)
				}
			}
		}
	})

	saveBtn := js.Global.Get("document").Call("getElementById", "save-btn")
	saveBtn.Call("addEventListener", "click", func() {
		w := &ExportWriter{}
		blk, _ := t.GetBlockAt(0, 0)
		builder.WriteBlock(w, blk)
		w.Export()
	})

	loadInput := js.Global.Get("document").Call("getElementById", "load-input")
	loadInput.Call("addEventListener", "change", func() {
		file := loadInput.Get("files").Get("0")
		if file == nil {
			return
		}

		go (func() {
			r := &ExportReader{}
			r.File = file
			blk := handler.ReadBlock(r)
			t.SetBlock(blk, 0)
			t.Draw()
		})()
	})

	return t
}
