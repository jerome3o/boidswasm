package main

import (
	"fmt"
	"math"
	"syscall/js"
)

type BoidsState struct {
	Boids [][]float64
}

func main() {
	fmt.Println("Boids Online")
	js.Global().Set("updateBoids", js.FuncOf(updateBoidsWrapper()))
	<-make(chan bool)
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
		boidsState.Boids[0][2] += math.Pi / 100
		return boidsState
	}, nil

}

func updateBoidsWrapper() func(this js.Value, args []js.Value) interface{} {

	boidsF, err := updateBoids()

	return func(this js.Value, args []js.Value) interface{} {
		if err := checkArgs(args); err != nil {
			return convertError(err)
		}

		boids := boidsF()

		if err != nil {
			return convertError(err)
		}

		return boidsOutputToJsFriendly(boids)
	}
}

func checkArgs(args []js.Value) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 args, found %v", len(args))
	}
	return nil
}

func convertError(e error) map[string]interface{} {
	return map[string]interface{}{
		"error": e.Error(),
	}
}

func boidsOutputToJsFriendly(boidsOutput BoidsState) map[string]interface{} {
	// TODO(j.swannack): with better understanding of go, you'll probably want
	// 	to revise this function to construct the output more efficiently

	output := make(map[string]interface{})

	boids := make([]interface{}, len(boidsOutput.Boids))
	for i, r := range boidsOutput.Boids {
		i_val := make([]interface{}, len(r))
		for j, c := range r {
			i_val[j] = c
		}
		boids[i] = i_val
	}

	output["boids"] = boids
	return output
}
