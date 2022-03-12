package main

import (
	"fmt"
	"syscall/js"
	"time"
)

type JsFunc func(this js.Value, args []js.Value) interface{}

var NFramesToAverage int = 10

func main() {
	fmt.Println("Boids Online")
	update, init := getWrappedBoidsFunctions()

	js.Global().Set("updateBoids", js.FuncOf(update))
	js.Global().Set("initBoids", js.FuncOf(init))
	<-make(chan bool)
}

func getWrappedBoidsFunctions() (JsFunc, JsFunc) {

	iFrame := 0
	var cumulativeFrameTime int64 = 0

	boidsF, init, err := getBoidsEngine()

	updateWrapped := func(this js.Value, args []js.Value) interface{} {

		iFrame += 1
		tStart := time.Now().UnixMilli()

		if err := checkUpdateArgs(args); err != nil {
			return convertError(err)
		}

		boids := boidsF(jsUpdateRequestToGo(args[0]))

		if err != nil {
			return convertError(err)
		}

		output := boidsOutputToJsFriendly(boids)

		// TODO(j.swannack): Could be nice patternf or this using defer?
		cumulativeFrameTime += time.Now().UnixMilli() - tStart
		if (iFrame % NFramesToAverage) == 0 {
			fmt.Printf("Average go calculation time: %vms\n", float64(cumulativeFrameTime)/float64(NFramesToAverage))
			cumulativeFrameTime = 0
		}

		return output
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
	settings := make(map[string]interface{})
	for k := range boidsOutput.Settings {
		settings[k] = boidsOutput.Settings[k]
	}
	output["settings"] = settings

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
	output := make(BoidSettings)

	for k := range defaultSettings {
		if !jsv.Get(k).IsUndefined() {
			output[k] = jsv.Get(k).Float()
		}
	}
	return output
}

func jsUpdateRequestToGo(jsv js.Value) BoidsUpdateRequest {
	output := BoidsUpdateRequest{}
	output.TimeStep = jsv.Get("timeStep").Float()
	output.Settings = jsSettingsToGo(jsv.Get("settings"))
	output.MouseX = jsv.Get("mouseX").Float()
	output.MouseY = jsv.Get("mouseY").Float()
	return output
}
