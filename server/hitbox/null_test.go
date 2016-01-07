package hitbox_test

import (
	"testing"

	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
)

func TestNull(t *testing.T) {
	hb := hitbox.NewNull()

	contour := hb.Contour(nil)
	if len(contour) != 0 {
		t.Fatal("Null hitbox contour must be empty")
	}
}
