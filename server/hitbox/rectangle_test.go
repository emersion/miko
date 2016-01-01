package hitbox_test

import (
	"testing"

	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

func TestRectangle(t *testing.T) {
	hb := hitbox.NewRectangle(&terrain.Position{10, 22}, 4, 6)

	contour := hb.Contour()

	if len(contour) != 4 {
		t.Fatal("Rectangle hitbox contour must contain 4 points")
	}

	if !contour[0].Equals(&terrain.Position{8, 19}) {
		t.Error("Invalid top-left corner position")
	}
	if !contour[1].Equals(&terrain.Position{12, 19}) {
		t.Error("Invalid top-right corner position")
	}
	if !contour[2].Equals(&terrain.Position{12, 25}) {
		t.Error("Invalid bottom-right corner position")
	}
	if !contour[3].Equals(&terrain.Position{8, 25}) {
		t.Error("Invalid bottom-left corner position")
	}
}
