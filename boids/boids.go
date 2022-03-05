package main

import (
	"fmt"
	"math"
	"math/rand"
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

func updateBoids() (func(t float64) BoidsState, func(h, w int) BoidsState, error) {
	isInit := false
	var boidsState BoidsState

	tTotal := 0.0

	var height float64
	var width float64

	dMax := 100.0
	vMax := 100.0

	sFactor := 1000.0
	cFactor := 10.0

	update := func(t float64) BoidsState {
		tTotal += t
		if !isInit {
			return boidsState
		}

		for i, boid := range boidsState.Boids {
			x, y, vx, vy := boid[0], boid[1], boid[2], boid[3]
			nearBoids := getNearBoids(x, y, dMax, i, boidsState.Boids)

			cax, cay := calculateCohesionDeltaV(x, y, vx, vy, vMax, nearBoids)
			sax, say := calculateSeparationDeltaV(x, y, vx, vy, vMax, nearBoids)
			// sepAX, sepAY := moveTowardNearestBoid(x, y, vx, vy, nearBoids)

			vx = vx + sFactor*sax + cFactor*cax
			vy = vy + sFactor*say + cFactor*cay

			s := getDist(0, 0, vx, vy)
			if s > vMax {
				// if s > 0 {
				vx *= vMax / s
				vy *= vMax / s
			}

			x += vx * t
			y += vy * t
			boidsState.Boids[i][0] = wrap(x, width)
			boidsState.Boids[i][1] = wrap(y, height)
			boidsState.Boids[i][2] = vx
			boidsState.Boids[i][3] = vy
		}

		return boidsState
	}

	init := func(w, h int) BoidsState {
		fmt.Printf("%v, %v\n", w, h)

		width = float64(w)
		height = float64(h)

		nrows, ncols := h/100, w/100
		boidsState = BoidsState{}

		boidsState.Boids = make([][]float64, nrows*ncols)

		for i := 0; i < ncols; i++ {
			for j := 0; j < nrows; j++ {
				boidsState.Boids[i*nrows+j] = []float64{
					width / 10 * rand.Float64(),
					height / 10 * rand.Float64(),
					(rand.Float64() - 0.5) * 200,
					(rand.Float64() - 0.5) * 200,
				}
			}
		}

		isInit = true
		return boidsState
	}

	return update, init, nil

}

func calculateSeparationDeltaV(x, y, vx, vy, maxV float64, boids [][]float64) (ax, ay float64) {

	if len(boids) == 0 {
		return 0.0, 0.0
	}

	ax, ay = 0.0, 0.0
	for _, b := range boids {

		dx := x - b[0]
		dy := y - b[1]

		dxSign := 1.0
		if dx < 0 {
			dxSign = -1.0
		}
		dySign := 1.0
		if dy < 0 {
			dySign = -1.0
		}

		ax += math.Min(maxV/math.Max(math.Abs(dx), 1), maxV) * dxSign
		ay += math.Min(maxV/math.Max(math.Abs(dy), 1), maxV) * dySign
	}

	// fmt.Println(ax, ay)

	return ax / float64(len(boids)), ay / float64(len(boids))
}

func calculateCohesionDeltaV(x, y, vx, vy, maxV float64, boids [][]float64) (ax, ay float64) {

	if len(boids) == 0 {
		return 0.0, 0.0
	}

	// TODO(j.swannack): account for wrap - might need w, h to be passed in?

	xCentre := 0.0
	yCentre := 0.0

	for _, b := range boids {
		xCentre += (b[0] - x)
		yCentre += (b[1] - y)
	}
	xCentre /= float64(len(boids))
	yCentre /= float64(len(boids))

	return xCentre, yCentre
}

func getNearBoids(x, y, dMax float64, iBoid int, boids [][]float64) [][]float64 {
	output := make([][]float64, 0)

	for i, b := range boids {
		if i == iBoid {
			continue
		}
		if getDist(x, y, b[0], b[1]) < dMax {
			// TODO(j.swannack): Check if boid in field of view
			output = append(output, b)
		}
	}
	return output
}

func getDist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
