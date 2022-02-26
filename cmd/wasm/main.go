package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Boids Online")
	js.Global().Set("updateBoids", js.FuncOf(updateBoidsWrapper()))
	<-make(chan bool)
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
