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

	var height float64
	var width float64

	update := func(t float64) BoidsState {
		if !isInit {
			return boidsState
		}

		fmt.Println(width)

		for i, boid := range boidsState.Boids {
			x, y, a, v := boid[0], boid[1], boid[2], boid[3]
			a += 1 * t
			vx, vy := math.Cos(a)*v*t, math.Sin(a)*v*t

			x += vx
			y += vy
			v = math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
			a = math.Atan(vy / vx)
			boidsState.Boids[i][0] = math.Mod(x, width)
			boidsState.Boids[i][1] = math.Mod(y, height)
			boidsState.Boids[i][2] = a
			boidsState.Boids[i][3] = v
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
					float64(i*j) / 100,
					1.0,
				}
			}
		}

		isInit = true
	}

	return update, init, nil

}
