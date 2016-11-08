package hitbox_test

import (
	"testing"

	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

func TestPoint(t *testing.T) {
	hb := hitbox.NewPoint()
	contour := hb.Contour(&terrain.Position{7, 22})

	if len(contour) != 1 {
		t.Fatal("Point hitbox contour must contain 1 point")
	}

	if !contour[0].Equals(&terrain.Position{7, 22}) {
		t.Error("Invalid contour")
	}
}
