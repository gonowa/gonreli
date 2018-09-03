BINARY=go.wasm
BINARY_OPTIMIZED=goO.wasm

all:    deeps   build

build:
	GOOS=js GOARCH=wasm CGO_ENABLED=0 go build -v -o ${BINARY}

deeps:
	GOOS=js GOARCH=wasm CGO_ENABLED=0 go get -t -v
	go get -u github.com/gonowa/wasm-opt

optimized:
	wasm-opt ${BINARY} -O4 -d -o ${BINARY_OPTIMIZED}
