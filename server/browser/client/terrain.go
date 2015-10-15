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
	canvas *canvas.Canvas
}

func (t *Terrain) DrawPoint(x, y int) {
	ptType := t.Points[x][y]

	if int(ptType) == 0 {
		t.canvas.ClearRect(x * res, y * res, res, res)
	} else {
		t.canvas.FillRect(x * res, y * res, res, res)
	}
}
func (t *Terrain) DrawRegion(x, y, w, h int) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
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
	t.canvas.Element.Set("width", size * message.BLOCK_LEN * res)
	t.canvas.Element.Set("height", size * message.BLOCK_LEN * res)
}

func NewTerrain(el *js.Object) *Terrain {
	t := &Terrain{*terrain.New(),canvas.New(el)}
	t.Reset(2)
	return t
}
