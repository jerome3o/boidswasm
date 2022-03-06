const go = new Go();

let w = window.innerWidth > 0 ? window.innerWidth : screen.width;
let h = window.innerHeight > 0 ? window.innerHeight : screen.height;
let boidsInitialised = false
let lastMs = 0.0

let debugNextPrint = 1000
let debugNextPrintStep = 100000

WebAssembly.instantiateStreaming(fetch("boids.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
}).then(v => {
    console.log(initBoids(w, h));
    boidsInitialised = true  // TODO(j.swannack): actually check for successful init
});

function setup() {
    createCanvas(w, h)
    lastMs = millis()
}

function draw() {

    if (!boidsInitialised) {
        return
    }

    ms = millis()
    timeStep = (ms - lastMs) / 1000
    lastMs = ms

    boids = updateBoids(timeStep)
    background(255)
    boids.boids.map(v => drawBoid(...v))


    if (ms > debugNextPrint) {
        debugNextPrint += debugNextPrintStep
        console.log(boids)
    }
}

function drawBoid(x, y, vx, vy) {
    push()

    stroke("black")

    let a = Math.atan(vy/vx)
    if (vx < 0) {
        a = a + Math.PI
    }

    translate(x, y)
    rotate(a)

    line(5, 0, -5,  6)
    line(5, 0, -5, -6)

    pop()
}