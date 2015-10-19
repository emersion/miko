package client

import (
	"github.com/gopherjs/gopherjs/js"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko/server/message/handler"
	"git.emersion.fr/saucisse-royale/miko/server/browser/client"
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
	el.Set("width", size * message.BLOCK_LEN * res)
	el.Set("height", size * message.BLOCK_LEN * res)

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
		t.DrawRegion(fromX, fromY, lastX - fromX, lastY - fromY)
		t.Canvas.SetFillStyle("rgba(0,0,255,0.1)")
		t.Canvas.FillRect(fromX * res, fromY * res, (x - fromX) * res, (y - fromY) * res)
		lastX, lastY = x, y
	})
	el.Call("addEventListener", "mouseup", func(event *js.Object) {
		t.DrawRegion(fromX, fromY, lastX - fromX, lastY - fromY)
		pressing = false

		if fromX == lastX && fromY == lastY {
			x, y := fromX, fromY
			t.Points[x][y] = invertPoint(t.Points[x][y])
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
					t.Points[i][j] = invertPoint(t.Points[i][j])
					t.DrawPoint(i, j)
				}
			}
		}
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
