package main

import(
	"git.emersion.fr/saucisse-royale/miko/server/browser/mapper/client"
)

func main() {
	trn := client.NewTerrain()
	trn.Draw()
}
