package main

import (
	"fmt"
	"math"
)

type BoidsState struct {
	Boids [][]float64
}

func updateBoids() (func(t float64) BoidsState, func(h, w int), error) {
	isInit := false
	var boidsState BoidsState

	update := func(t float64) BoidsState {
		if !isInit {
			return boidsState
		}

		for _, boid := range boidsState.Boids {
			boid[2] += t * math.Pi
		}

		return boidsState
	}

	init := func(w, h int) {
		fmt.Printf("%v, %v\n", w, h)

		nrows, ncols := h/100, w/100
		boidsState = BoidsState{}

		boidsState.Boids = make([][]float64, nrows*ncols)

		for i := 0; i < ncols; i++ {
			for j := 0; j < nrows; j++ {
				boidsState.Boids[i*nrows+j] = []float64{
					float64(100*i) + 50,
					float64(100*j) + 50,
					float64(i*j) / 100,
					0.0,
				}
			}
		}

		isInit = true
	}

	return update, init, nil

}
