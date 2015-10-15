package client

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/emersion/go-js-canvas"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/terrain"
)

const res = 5

type Terrain struct {
	terrain.Terrain
	Canvas *canvas.Canvas
}

func (t *Terrain) DrawPoint(x, y int) {
	ptType := t.Points[x][y]

	if int(ptType) == 0 {
		t.Canvas.ClearRect(x * res, y * res, res, res)
	} else {
		t.Canvas.SetFillStyle("black")
		t.Canvas.FillRect(x * res, y * res, res, res)
	}
}

func (t *Terrain) DrawRegion(x, y, w, h int) {
	x1 := x
	y1 := y
	x2 := x + w
	y2 := y + h

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			t.DrawPoint(i, j)
		}
	}
}

func (t *Terrain) Draw() {
	for i := range t.Points {
		for j := range t.Points[i] {
			t.DrawPoint(i, j)
		}
	}
}

func (t *Terrain) SetBlock(blk *message.Block) {
	t.Terrain.SetBlock(blk)
	t.Draw()
}

func (t *Terrain) Reset(size int) {
	t.Terrain.Reset(size)
	t.Canvas.Element.Set("width", size * message.BLOCK_LEN * res)
	t.Canvas.Element.Set("height", size * message.BLOCK_LEN * res)
}

func NewTerrain(el *js.Object) *Terrain {
	t := &Terrain{*terrain.New(),canvas.New(el)}
	t.Reset(2)
	return t
}
