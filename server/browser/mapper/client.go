package main

import(
	"git.emersion.fr/saucisse-royale/miko/server/browser/client"
)

func main() {
	trn := client.NewTerrain()
	trn.Draw()
}
