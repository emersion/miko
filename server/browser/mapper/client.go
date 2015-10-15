package main

import(
	"github.com/gopherjs/gopherjs/js"

	"git.emersion.fr/saucisse-royale/miko/server/browser/mapper/terrain"
)

func main() {
	canvas := js.Global.Get("document").Call("getElementById", "canvas")
	trn := client.NewTerrain(canvas)
	trn.Draw()
}
