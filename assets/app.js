const go = new Go();

let w = window.innerWidth > 0 ? window.innerWidth : screen.width;
let h = window.innerHeight > 0 ? window.innerHeight : screen.height;
let boidsInitialised = false
let lastMs = 0.0

let debugNextPrint = 1000
let debugNextPrintStep = 100000


let sliderSpacing = 35
let sliderWidth = 200

let sliderSpec = [
    {
        title: "Neighbourhood Radius",
        settingKey: "distMax",
        min: 0,
        max: 600,
        default: 50.0
    },
    {
        title: "Maximum Velocity",
        settingKey: "velocityMax",
        min: 0,
        max: 1000,
        default: 200.0
    },
    {
        title: "Separation",
        settingKey: "separationFactor",
        min: 0,
        max: 50,
        default: 3.0
    },
    {
        title: "Cohesion",
        settingKey: "cohesionFactor",
        min: 0,
        max: 10,
        default: 1.0
    },
    {
        title: "Alignment",
        settingKey: "alignmentFactor",
        min: 0,
        max: 10,
        default: 3.0
    },
    {
        title: "Random Movement",
        settingKey: "randomFactor",
        min: 0,
        max: 100,
        default: 1.0
    },
]
let sliders = {}

WebAssembly.instantiateStreaming(fetch("boids.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
}).then(v => {
    console.log(initBoids(w, h));
    boidsInitialised = true  // TODO(j.swannack): actually check for successful init
});

function setup() {
    createCanvas(w, h)

    for (let i = 0; i < sliderSpec.length; i++) {
        let spec = sliderSpec[i]
        let slider = createSlider(spec.min, spec.max, spec.default, 1)
        slider.position(20, sliderSpacing*(i+1))
        slider.style('width', `${sliderWidth}px`);
        sliders[spec.settingKey] = slider
    }

    lastMs = millis()
}

function createSlider(title, index){

    let spacing = 20

    let slider = createSlider(0, 800, 50, 1)
    slider.position(20, sliderSpacing*(index+1))
    slider.style('width', '200px');
    text(title, slider.x * 2 + slider.width, sliderSpacing*(index+1));
    return slider
}

function draw() {

    if (!boidsInitialised) {
        return
    }

    ms = millis()
    timeStep = (ms - lastMs) / 1000
    lastMs = ms

    let settings = {}
    for(var key in sliders) {
        settings[key] = sliders[key].value()
    }

    boids = updateBoids({
        timeStep,
        settings,
    })
    background(255)
    boids.boids.map(v => drawBoid(...v))

    boids.debugBoids.map(db => drawDebugBoid(
        boids.boids[db.index],
        db.neighbours.map(i => boids.boids[i]),
        boids.settings,
    ))
    
    if (ms > debugNextPrint) {
        debugNextPrint += debugNextPrintStep
        console.log(boids)
    }

    drawSliderText()
}

function drawSliderText() {
    push()
    stroke("black")
    strokeWeight(1)
    for (let i = 0; i < sliderSpec.length; i++) {
        let spec = sliderSpec[i]
        let slider = sliders[spec.settingKey]
        text(`${spec.title}: ${slider.value()}`, slider.x + slider.width + 10, sliderSpacing*(i+1) + 15);
    }
    pop()
}

function drawDebugBoid(boid, neighbours, settings) {
    push()
    stroke("purple")
    fill(0,0,0,0)
    strokeWeight(1)
    circle(boid[0], boid[1], 20)
    circle(boid[0], boid[1], settings.distMax*2)

    if (boid[0] < settings.distMax) {
        circle(boid[0] + settings.width, boid[1], settings.distMax*2)
    }

    if (boid[0] + settings.distMax > settings.width) {
        circle(boid[0] - settings.width, boid[1], settings.distMax*2)
    }

    if (boid[1] < settings.distMax) {
        circle(boid[0], boid[1] + settings.height, settings.distMax*2)
    }

    if (boid[1] + settings.distMax > settings.height) {
        circle(boid[0], boid[1] - settings.height, settings.distMax*2)
    }

    strokeWeight(1)
    neighbours.map(b => circle(b[0], b[1], 20))

    pop()
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

    line(3, 0, -3,  3)
    line(3, 0, -3, -3)

    pop()
}