package main

import "math"

type BoidsState struct {
	Boids [][]float64
}

func updateBoids() (func() BoidsState, error) {
	boidsState := BoidsState{
		Boids: [][]float64{
			{50, 50, 0, 10},
			{100, 50, math.Pi / 2, 10},
			{150, 50, math.Pi, 10},
			{200, 50, 3 * math.Pi / 2, 10},
			{250, 50, 0, 10},
		},
	}

	return func() BoidsState {
		boidsState.Boids[0][2] -= math.Pi / 100
		return boidsState
	}, nil

}
