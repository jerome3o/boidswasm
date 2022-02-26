# Boidswa(s/r)m

Playing around trying to get boids going in the browser using go for computation and p5js for rendering

# Build WASM

```
./build.sh
```

# Run server

## With just go

```
go run server/main.go
```

This is handy while doing dev to rebuild and serve

```
./build.sh && go run server/main.go
```

## With browser-sync for live reloading

```bash
cd assets
browser-sync start -s -f . --no-notify --host 0.0.0.0 --port 9000
```

# VSCode

for your `.vscode/settings.json` file:

```json
{
    "go.testEnvVars": {
        "GOOS": "js",
        "GOARCH": "wasm"
    },
    "go.toolsEnvVars": {
        "GOOS": "js",
        "GOARCH": "wasm"
    }
}
```