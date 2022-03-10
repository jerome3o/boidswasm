package main

import (
	"fmt"
	"math"
	"math/rand"
)

type BoidSettings map[string]float64

type DebugBoid struct {
	Index      int
	Neighbours []int
}

type BoidsState struct {
	Boids      [][]float64
	Settings   BoidSettings
	DebugBoids []DebugBoid
}

type BoidsUpdateRequest struct {
	TimeStep float64
	Settings BoidSettings
	MouseX   float64
	MouseY   float64
}

var defaultSettings BoidSettings = map[string]float64{
	"distMax":          100.0,
	"velocityMax":      300.0,
	"separationFactor": 10.0,
	"cohesionFactor":   1.0,
	"alignmentFactor":  1.0,
	"randomFactor":     1.0,
	"fearFactor":       1.0,
	"width":            1000.0,
	"height":           1000.0,
}

func wrap(x, bound float64) float64 {
	for x < 0 {
		x += bound
	}
	return math.Mod(x, bound)
}

func updateBoids() (func(update BoidsUpdateRequest) BoidsState, func(h, w int) BoidsState, error) {
	isInit := false
	var boidsState BoidsState

	tTotal := 0.0

	update := func(updateReq BoidsUpdateRequest) BoidsState {

		for k, v := range updateReq.Settings {
			boidsState.Settings[k] = v
		}

		t := updateReq.TimeStep
		dMax := boidsState.Settings["distMax"]
		vMax := boidsState.Settings["velocityMax"]
		sFactor := boidsState.Settings["separationFactor"]
		cFactor := boidsState.Settings["cohesionFactor"]
		aFactor := boidsState.Settings["alignmentFactor"]
		rFactor := boidsState.Settings["randomFactor"]
		fFactor := boidsState.Settings["fearFactor"]
		height := boidsState.Settings["height"]
		width := boidsState.Settings["width"]

		mouseX := updateReq.MouseX
		mouseY := updateReq.MouseY

		tTotal += t
		if !isInit {
			return boidsState
		}

		newBoids := make([][]float64, len(boidsState.Boids))

		for i, boid := range boidsState.Boids {
			x, y := boid[0], boid[1]
			nearBoidIndices, nearBoids := getNearBoids(x, y, width, height, dMax, i, boidsState.Boids)

			for ii, debugBoid := range boidsState.DebugBoids {
				if i == debugBoid.Index {
					boidsState.DebugBoids[ii].Neighbours = nearBoidIndices
				}
			}

			// TODO(j.swannack): these guys need a good debug - they get overly grouped
			cax, cay := calculateCohesionDeltaV(x, y, width, height, vMax, nearBoids)
			sax, say := calculateSeparationDeltaV(x, y, width, height, dMax, nearBoids)
			aax, aay := calculateAlignmentDeltaV(x, y, width, height, vMax, nearBoids)
			fax, fay := calculateFearDeltaV(x, y, mouseX, mouseY, width, height, vMax, nearBoids)

			// TODO(j.swannack): Think more about this
			rax, ray := math.Sin(float64(i)+tTotal), math.Cos(float64(i)+tTotal)

			vx := sFactor*sax + cFactor*cax + aFactor*aax + rFactor*rax + fFactor*fax
			vy := sFactor*say + cFactor*cay + aFactor*aay + rFactor*ray + fFactor*fay

			s := getDist(0, 0, vx, vy)
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

		nrows, ncols := h/80, w/80
		boidsState = BoidsState{}

		boidsState.Boids = make([][]float64, nrows*ncols)
		boidsState.Settings = BoidSettings{
			"distMax":          50.0,
			"velocityMax":      200.0,
			"separationFactor": 3.0,
			"cohesionFactor":   1.0,
			"alignmentFactor":  3.0,
			"randomFactor":     1.0,
			"fearFactor":       1.0,
			"width":            float64(w),
			"height":           float64(h),
		}

		for i := 0; i < ncols; i++ {
			for j := 0; j < nrows; j++ {
				boidsState.Boids[i*nrows+j] = []float64{
					boidsState.Settings["width"] * rand.Float64(),
					boidsState.Settings["height"] * rand.Float64(),
					(rand.Float64() - 0.5) * 200,
					(rand.Float64() - 0.5) * 200,
				}
			}
		}

		boidsState.DebugBoids = []DebugBoid{
			{
				Index:      0,
				Neighbours: []int{},
			},
		}

		isInit = true
		return boidsState
	}

	return update, init, nil

}

func calculateSeparationDeltaV(x, y, w, h, distMax float64, boids [][]float64) (ax, ay float64) {

	if len(boids) == 0 {
		return 0.0, 0.0
	}

	ax, ay = 0.0, 0.0
	for _, b := range boids {

		dx := getWrappedDist1d(b[0], x, w)
		dy := getWrappedDist1d(b[1], y, h)

		dxSign := 1.0
		if dx < 0 {
			dxSign = -1.0
		}
		dySign := 1.0
		if dy < 0 {
			dySign = -1.0
		}

		ax += math.Min(distMax/math.Max(math.Abs(dx), 1), 2*distMax) * dxSign
		ay += math.Min(distMax/math.Max(math.Abs(dy), 1), 2*distMax) * dySign
	}

	return ax, ay
}

func calculateCohesionDeltaV(x, y, w, h, maxV float64, boids [][]float64) (ax, ay float64) {

	xCentre := 0.0
	yCentre := 0.0

	if len(boids) == 0 {
		return xCentre, yCentre
	}

	for _, b := range boids {
		xCentre += getWrappedDist1d(x, b[0], w)
		yCentre += getWrappedDist1d(y, b[1], h)
	}
	xCentre /= float64(len(boids))
	yCentre /= float64(len(boids))

	return xCentre, yCentre
}

func calculateAlignmentDeltaV(x, y, w, h, maxV float64, boids [][]float64) (ax, ay float64) {

	vxAv := 0.0
	vyAv := 0.0

	if len(boids) == 0 {
		return vxAv, vyAv
	}

	for _, b := range boids {
		vxAv += b[2]
		vyAv += b[3]
	}
	vxAv /= float64(len(boids)) * 10
	vyAv /= float64(len(boids)) * 10

	return vxAv, vyAv
}

func calculateFearDeltaV(x, y, mouseX, mouseY, width, height, vMax float64, boids [][]float64) (ax, ay float64) {

	dist := getWrappedDist(x, y, mouseX, mouseY, width, height)

	if dist > 100 {
		return 0.0, 0.0
	}

	ax = getWrappedDist1d(mouseX, x, width)
	ay = getWrappedDist1d(mouseY, y, height)

	return ax, ay
}

func getNearBoids(x, y, w, h, dMax float64, iBoid int, boids [][]float64) ([]int, [][]float64) {
	output := make([][]float64, 0)
	indices := make([]int, 0)

	for i, b := range boids {
		if i == iBoid {
			continue
		}
		if getWrappedDist(x, y, b[0], b[1], w, h) < dMax {
			// TODO(j.swannack): Check if boid in field of view
			output = append(output, b)
			indices = append(indices, i)
		}
	}
	return indices, output
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

	dist := v2 - v1
	if math.Abs(dist) > bound/2 {
		if dist < 0 {
			return dist + bound
		} else {
			return dist - bound
		}
	}
	return dist
}
