const go = new Go();
WebAssembly.instantiateStreaming(fetch("boids.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
}).then(v => {
    console.log(updateBoids());
});

