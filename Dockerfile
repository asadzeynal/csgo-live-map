FROM golang:1.19.5-alpine3.17
WORKDIR /app
COPY . .

RUN GOOS=js GOARCH=wasm go build -o static/main.wasm cmd/wasm/*

EXPOSE 8080

ENTRYPOINT ["go", "run", "cmd/webserver/main.go"]
