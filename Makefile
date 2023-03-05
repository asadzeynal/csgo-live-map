build:
	GOOS=js GOARCH=wasm go build -o static/main.wasm cmd/wasm/*
run:
	go run ./cmd/webserver/main.go
