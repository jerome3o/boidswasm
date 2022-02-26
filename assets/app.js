const go = new Go();
WebAssembly.instantiateStreaming(fetch("boids.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
}).then(v => {
    console.log(updateBoids());
});

let w = window.innerWidth > 0 ? window.innerWidth : screen.width;
let h = window.innerHeight > 0 ? window.innerHeight : screen.height;

function setup() {
    createCanvas(w, h)
}

function draw() {
    boids = updateBoids()
    background(0)
    boids.boids.map(v => drawBoid(...v))
}

function drawBoid(x, y, a, v) {
    push()

    stroke("white")

    translate(x, y)
    rotate(a)

    line(0, 5, 3, -5)
    line(0, 5, -3, -5)

    pop()
}