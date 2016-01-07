package hitbox_test

import (
	"testing"

	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

func TestCircle(t *testing.T) {
	hb := hitbox.NewCircle(10)
	contour := hb.Contour(&terrain.Position{10, 22})

	if len(contour) != 4 {
		t.Fatal("Circle hitbox contour must contain 4 points")
	}

	expected := []*terrain.Position{
		&terrain.Position{20, 22},
		&terrain.Position{10, 32},
		&terrain.Position{0, 22},
		&terrain.Position{10, 12},
	}

	// Compare rounded values, because of precision errors
	for i, pt := range contour {
		if !pt.ToMessage().Equals(expected[i].ToMessage()) {
			t.Error("Invalid contour point at index", i, "expected", expected[i], "but got", pt)
		}
	}
}
