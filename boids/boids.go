package main

import (
	"fmt"
	"math"
)

type BoidsState struct {
	Boids [][]float64
}

func wrap(x, bound float64) float64 {
	for x < 0 {
		x += bound
	}
	return math.Mod(x, bound)
}

func updateBoids() (func(t float64) BoidsState, func(h, w int), error) {
	isInit := false
	var boidsState BoidsState

	tTotal := 0.0

	var height float64
	var width float64

	update := func(t float64) BoidsState {
		tTotal += t
		if !isInit {
			return boidsState
		}

		for i, boid := range boidsState.Boids {
			x, y, _, vy := boid[0], boid[1], boid[2], boid[3]

			vx := math.Sin(tTotal*1+float64(i)) * 100
			// vy := math.Cos(tTotal*1+float64(i)) * 100
			x += vx * t
			y += vy * t
			boidsState.Boids[i][0] = wrap(x, width)
			boidsState.Boids[i][1] = wrap(y, height)
			boidsState.Boids[i][2] = vx
			boidsState.Boids[i][3] = vy
		}

		return boidsState
	}

	init := func(w, h int) {
		fmt.Printf("%v, %v\n", w, h)

		width = float64(w)
		height = float64(h)

		nrows, ncols := h/100, w/100
		boidsState = BoidsState{}

		boidsState.Boids = make([][]float64, nrows*ncols)

		for i := 0; i < ncols; i++ {
			for j := 0; j < nrows; j++ {
				boidsState.Boids[i*nrows+j] = []float64{
					float64(100*i) + 50,
					float64(100*j) + 50,
					0.0,
					100.0,
				}
			}
		}

		isInit = true
	}

	return update, init, nil

}
