package terrain_test

import (
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

func ExampleGetRouteBetween_point() {
	from := &terrain.Position{0, 0}
	to := &terrain.Position{0, 0}

	route := terrain.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[0 0]]
}

func ExampleGetRouteBetween_horizontal() {
	from := &terrain.Position{-2, 0}
	to := &terrain.Position{4, 0}

	route := terrain.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[-2 0] [-1 0] [0 0] [1 0] [2 0] [3 0] [4 0]]
}

func ExampleGetRouteBetween_diagonal() {
	from := &terrain.Position{-2, 2}
	to := &terrain.Position{2, -2}

	route := terrain.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[-2 2] [-1 1] [0 0] [1 -1] [2 -2]]
}

func ExampleGetRouteBetween_other() {
	from := &terrain.Position{-1, 5}
	to := &terrain.Position{5, 0}

	route := terrain.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[-1 5] [0 4] [1 3] [2 3] [3 2] [4 1] [5 0]]
}
