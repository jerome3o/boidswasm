package main

import (
	"fmt"
	"math"
)

type BoidsState struct {
	Boids [][]float64
}

func updateBoids() (func() BoidsState, func(h, w int), error) {
	isInit := false
	var boidsState BoidsState

	update := func() BoidsState {
		if !isInit {
			return boidsState
		}

		boidsState.Boids[0][2] += math.Pi / 10
		return boidsState
	}

	init := func(w, h int) {
		fmt.Printf("%v, %v\n", w, h)
		boidsState = BoidsState{
			Boids: [][]float64{
				{50, 50, 0, 10},
				{100, 50, math.Pi / 2, 10},
				{150, 50, math.Pi, 10},
				{200, 50, 3 * math.Pi / 2, 10},
				{250, 50, 0, 10},
			},
		}
		isInit = true
	}

	return update, init, nil

}
