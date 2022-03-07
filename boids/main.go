package main

import (
	"fmt"
	"syscall/js"
)

type JsFunc func(this js.Value, args []js.Value) interface{}

func main() {
	fmt.Println("Boids Online")
	update, init := getWrappedBoidsFunctions()

	js.Global().Set("updateBoids", js.FuncOf(update))
	js.Global().Set("initBoids", js.FuncOf(init))
	<-make(chan bool)
}

func getWrappedBoidsFunctions() (JsFunc, JsFunc) {

	boidsF, init, err := updateBoids()

	updateWrapped := func(this js.Value, args []js.Value) interface{} {
		if err := checkUpdateArgs(args); err != nil {
			return convertError(err)
		}

		boids := boidsF(jsUpdateRequestToGo(args[0]))

		if err != nil {
			return convertError(err)
		}

		return boidsOutputToJsFriendly(boids)
	}

	initWrapped := func(this js.Value, args []js.Value) interface{} {
		if err := checkInitArgs(args); err != nil {
			return convertError(err)
		}

		boids := init(
			args[0].Int(),
			args[1].Int(),
		)

		if err != nil {
			return convertError(err)
		}

		return boidsOutputToJsFriendly(boids)
	}

	return updateWrapped, initWrapped
}

func checkUpdateArgs(args []js.Value) error {
	// TODO(j.swannack): Check types
	if len(args) != 1 {
		return fmt.Errorf("expected 1 args, found %v", len(args))
	}
	return nil
}

func checkInitArgs(args []js.Value) error {
	// TODO(j.swannack): This boiler is surely abstracted/abstractable
	// TODO(j.swannack): Check types
	if len(args) != 2 {
		return fmt.Errorf("expected 2 args, found %v", len(args))
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
	output["settings"] = map[string]interface{}{
		"distMax":          boidsOutput.Settings.DistMax,
		"velocityMax":      boidsOutput.Settings.VelocityMax,
		"separationFactor": boidsOutput.Settings.SeparationFactor,
		"cohesionFactor":   boidsOutput.Settings.CohesionFactor,
		"alignmentFactor":  boidsOutput.Settings.AlignmentFactor,
		"width":            boidsOutput.Settings.Width,
		"height":           boidsOutput.Settings.Height,
	}

	debugBoids := make([]interface{}, len(boidsOutput.DebugBoids))

	for i, debugBoid := range boidsOutput.DebugBoids {

		neighbours := make([]interface{}, len(debugBoid.Neighbours))
		for ii, v := range debugBoid.Neighbours {
			neighbours[ii] = v
		}

		debugBoids[i] = map[string]interface{}{
			"index":      debugBoid.Index,
			"neighbours": neighbours,
		}
	}
	output["debugBoids"] = debugBoids

	return output
}

func jsSettingsToGo(jsv js.Value) BoidSettings {
	if jsv.IsUndefined() {
		return BoidSettings{}
	}

	return BoidSettings{
		DistMax:          jsv.Get("distMax").Float(),
		VelocityMax:      jsv.Get("velocityMax").Float(),
		SeparationFactor: jsv.Get("separationFactor").Float(),
		CohesionFactor:   jsv.Get("cohesionFactor").Float(),
		AlignmentFactor:  jsv.Get("alignmentFactor").Float(),
		Height:           jsv.Get("height").Float(),
		Width:            jsv.Get("width").Float(),
	}
}

func jsUpdateRequestToGo(jsv js.Value) BoidsUpdateRequest {
	output := BoidsUpdateRequest{}
	output.TimeStep = jsv.Get("timeStep").Float()
	output.Settings = jsSettingsToGo(jsv.Get("settings"))
	return output
}
