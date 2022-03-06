package main

import (
	"fmt"
	"math"
	"math/rand"
)

type BoidSettings struct {
	DistMax          float64
	VelocityMax      float64
	SeparationFactor float64
	CohesionFactor   float64
	AlignmentFactor  float64
	Height           float64
	Width            float64
}

type DebugBoid struct {
	Index      int
	Neighbours []int
}

type BoidsState struct {
	Boids    [][]float64
	Settings BoidSettings
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

	update := func(t float64) BoidsState {

		dMax := boidsState.Settings.DistMax
		vMax := boidsState.Settings.VelocityMax
		sFactor := boidsState.Settings.SeparationFactor
		cFactor := boidsState.Settings.CohesionFactor
		aFactor := boidsState.Settings.AlignmentFactor
		height := boidsState.Settings.Height
		width := boidsState.Settings.Width

		tTotal += t
		if !isInit {
			return boidsState
		}

		newBoids := make([][]float64, len(boidsState.Boids))

		for i, boid := range boidsState.Boids {
			x, y, vx, vy := boid[0], boid[1], boid[2], boid[3]
			nearBoids := getNearBoids(x, y, width, height, dMax, i, boidsState.Boids)

			cax, cay := calculateCohesionDeltaV(x, y, vx, vy, vMax, nearBoids)
			sax, say := calculateSeparationDeltaV(x, y, vx, vy, vMax, nearBoids)
			aax, aay := calculateAlignmentDeltaV(x, y, vx, vy, vMax, nearBoids)

			vx += sFactor*sax + cFactor*cax + aFactor*aax
			vy += sFactor*say + cFactor*cay + aFactor*aay

			s := getDist(0, 0, vx, vy)
			// if s > vMax {
			if s > 0 {
				vx *= vMax / s
				vy *= vMax / s
			}

			x += vx * t
			y += vy * t
			newBoids[i] = make([]float64, 4)
			newBoids[i][0] = wrap(x, width)
			newBoids[i][1] = wrap(y, height)
			newBoids[i][2] = vx
			newBoids[i][3] = vy
		}

		boidsState.Boids = newBoids
		return boidsState
	}

	init := func(w, h int) BoidsState {
		fmt.Printf("%v, %v\n", w, h)

		nrows, ncols := h/100, w/100
		boidsState = BoidsState{}

		boidsState.Boids = make([][]float64, nrows*ncols)
		boidsState.Settings = BoidSettings{
			DistMax:          100.0,
			VelocityMax:      300.0,
			SeparationFactor: 10.0,
			CohesionFactor:   1.0,
			AlignmentFactor:  1.0,
			Width:            float64(w),
			Height:           float64(h),
		}

		for i := 0; i < ncols; i++ {
			for j := 0; j < nrows; j++ {
				boidsState.Boids[i*nrows+j] = []float64{
					boidsState.Settings.Width * rand.Float64(),
					boidsState.Settings.Height * rand.Float64(),
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

		ax += math.Min(maxV/math.Max(math.Abs(dx), 1), 2*maxV) * dxSign
		ay += math.Min(maxV/math.Max(math.Abs(dy), 1), 2*maxV) * dySign
	}

	// fmt.Println(ax, ay)

	return ax / float64(len(boids)), ay / float64(len(boids))
}

func calculateCohesionDeltaV(x, y, vx, vy, maxV float64, boids [][]float64) (ax, ay float64) {
	// TODO(j.swannack): account for wrap - might need w, h to be passed in?

	xCentre := 0.0
	yCentre := 0.0

	if len(boids) == 0 {
		return xCentre, yCentre
	}

	for _, b := range boids {
		xCentre += (b[0] - x)
		yCentre += (b[1] - y)
	}
	xCentre /= float64(len(boids))
	yCentre /= float64(len(boids))

	return xCentre, yCentre
}

func calculateAlignmentDeltaV(x, y, vx, vy, maxV float64, boids [][]float64) (ax, ay float64) {

	vxAv := 0.0
	vyAv := 0.0

	if len(boids) == 0 {
		return vxAv, vyAv
	}

	for _, b := range boids {
		vxAv += b[2]
		vyAv += b[3]
	}
	vxAv /= float64(len(boids))
	vyAv /= float64(len(boids))

	return vxAv, vyAv
}

func getNearBoids(x, y, w, h, dMax float64, iBoid int, boids [][]float64) [][]float64 {
	output := make([][]float64, 0)

	for i, b := range boids {
		if i == iBoid {
			continue
		}
		if getWrappedDist(x, y, b[0], b[1], w, h) < dMax {
			// TODO(j.swannack): Check if boid in field of view
			output = append(output, b)
		}
	}
	return output
}

func getDist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func getWrappedDist(x1, y1, x2, y2, w, h float64) float64 {
	return math.Sqrt(
		math.Pow(getWrappedDist1d(x1, x2, w), 2) + math.Pow(getWrappedDist1d(y1, y2, h), 2),
	)
}

func getWrappedDist1d(v1, v2, bound float64) float64 {
	// TODO(j.swannack): Calculate distance that accounts for screen wrap
	// return wrap(v2-v1+bound/2.0, bound) - bound
	return v2 - v1
}
