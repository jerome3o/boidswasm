const go = new Go();

let w = window.innerWidth > 0 ? window.innerWidth : screen.width;
let h = window.innerHeight > 0 ? window.innerHeight : screen.height;
let boidsInitialised = false

WebAssembly.instantiateStreaming(fetch("boids.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
}).then(v => {
    console.log(initBoids(w, h));
    console.log(updateBoids());
    boidsInitialised = true  // TODO(j.swannack): actually check for successful init
});

function setup() {
    createCanvas(w, h)
}

function draw() {

    if (!boidsInitialised) {
        return
    }

    ms = millis()
    s = ms / 1000.0

    boids = updateBoids(s)
    background(255)
    boids.boids.map(v => drawBoid(...v))
}

function drawBoid(x, y, a, v) {
    push()

    stroke("black")

    translate(x, y)
    rotate(a)

    line(0, 5, 3, -5)
    line(0, 5, -3, -5)

    pop()
}