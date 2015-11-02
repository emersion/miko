package entity_test

import (
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

func ExampleGetRouteBetween_horizontal() {
	from := &entity.Position{-2, 0}
	to := &entity.Position{4, 0}

	route := entity.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[-2 0] [-1 0] [0 0] [1 0] [2 0] [3 0] [4 0]]
}

func ExampleGetRouteBetween_diagonal() {
	from := &entity.Position{-2, 2}
	to := &entity.Position{2, -2}

	route := entity.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[-2 2] [-1 1] [0 0] [1 -1] [2 -2]]
}

func ExampleGetRouteBetween_other() {
	from := &entity.Position{-1, 5}
	to := &entity.Position{5, 0}

	route := entity.GetRouteBetween(from, to)

	fmt.Println(route)
	// Output: [[-1 5] [0 4] [1 3] [2 3] [3 2] [4 1] [5 0]]
}
