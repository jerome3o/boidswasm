package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
)

var NBoids int = 1000

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

type BoidsEngine struct {
	boidsState BoidsState
	nextBoids  [][]float64
	isInit     bool
	tTotal     float64
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

func (e *BoidsEngine) Update(updateReq BoidsUpdateRequest) BoidsState {

	for k, v := range updateReq.Settings {
		e.boidsState.Settings[k] = v
	}

	t := updateReq.TimeStep
	dMax := e.boidsState.Settings["distMax"]
	vMax := e.boidsState.Settings["velocityMax"]
	sFactor := e.boidsState.Settings["separationFactor"]
	cFactor := e.boidsState.Settings["cohesionFactor"]
	aFactor := e.boidsState.Settings["alignmentFactor"]
	rFactor := e.boidsState.Settings["randomFactor"]
	fFactor := e.boidsState.Settings["fearFactor"]
	height := e.boidsState.Settings["height"]
	width := e.boidsState.Settings["width"]

	mouseX := updateReq.MouseX
	mouseY := updateReq.MouseY

	e.tTotal += t
	if !e.isInit {
		return e.boidsState
	}

	for i, boid := range e.boidsState.Boids {
		x, y, vx, vy := boid[0], boid[1], boid[2], boid[3]
		nearBoidIndices, nearBoids := getNearBoids(x, y, width, height, dMax, i, e.boidsState.Boids)

		for ii, debugBoid := range e.boidsState.DebugBoids {
			if i == debugBoid.Index {
				e.boidsState.DebugBoids[ii].Neighbours = nearBoidIndices
			}
		}

		// TODO(j.swannack): these guys need a good debug - they get overly grouped
		cax, cay := calculateCohesionDeltaV(x, y, width, height, vMax, nearBoids)
		sax, say := calculateSeparationDeltaV(x, y, width, height, dMax, nearBoids)
		aax, aay := calculateAlignmentDeltaV(x, y, width, height, vMax, nearBoids)
		fax, fay := calculateFearDeltaV(x, y, mouseX, mouseY, width, height, vMax, nearBoids)

		// TODO(j.swannack): Think more about this
		rax, ray := math.Sin(float64(i)+e.tTotal), math.Cos(float64(i)+e.tTotal)

		vx = sFactor*sax + cFactor*cax + aFactor*aax + rFactor*rax + fFactor*fax + 0.25*vx
		vy = sFactor*say + cFactor*cay + aFactor*aay + rFactor*ray + fFactor*fay + 0.25*vy

		s := getDist(0, 0, vx, vy)
		if s > 0 {
			vx *= vMax / s
			vy *= vMax / s
		}

		x += vx * t
		y += vy * t
		e.nextBoids[i][0] = wrap(x, width)
		e.nextBoids[i][1] = wrap(y, height)
		e.nextBoids[i][2] = vx
		e.nextBoids[i][3] = vy
	}

	for i, b := range e.boidsState.Boids {
		for ib := range b {
			e.boidsState.Boids[i][ib] = e.nextBoids[i][ib]
		}
	}
	return e.boidsState
}

func (e *BoidsEngine) Init(w, h int) BoidsState {
	fmt.Printf("%v, %v\n", w, h)
	fmt.Println("The number of CPU Cores:", runtime.NumCPU())

	e.boidsState = BoidsState{}

	e.nextBoids = make([][]float64, 2000)
	for i := range e.nextBoids {
		e.nextBoids[i] = make([]float64, 4)
	}

	e.boidsState.Boids = make([][]float64, NBoids)
	e.boidsState.Settings = BoidSettings{
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

	for i := 0; i < NBoids; i++ {
		e.boidsState.Boids[i] = []float64{
			e.boidsState.Settings["width"] * rand.Float64(),
			e.boidsState.Settings["height"] * rand.Float64(),
			(rand.Float64() - 0.5) * 200,
			(rand.Float64() - 0.5) * 200,
		}
	}

	e.boidsState.DebugBoids = []DebugBoid{
		{
			Index:      0,
			Neighbours: []int{},
		},
	}

	e.isInit = true
	return e.boidsState
}

func getBoidsEngine() (func(update BoidsUpdateRequest) BoidsState, func(h, w int) BoidsState, error) {
	isInit := false
	var boidsState BoidsState
	var nextBoids [][]float64

	tTotal := 0.0

	return update, init, nil

}

func wrap(x, bound float64) float64 {
	for x < 0 {
		x += bound
	}
	return math.Mod(x, bound)
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
