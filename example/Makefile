BINARY=go.wasm

all:    deeps   build

build:
	GOOS=js GOARCH=wasm CGO_ENABLED=0 go build -v -o ${BINARY}

deeps:
	GOOS=js GOARCH=wasm CGO_ENABLED=0 go get -t -v