package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

}

func compile() error {
	cmd := exec.Command("go", "test", "-c", "-o", "go.wasm", "github.com/gonowa/gonreli/test")
	env := os.Environ()
	env = append(env, fmt.Sprintf("CGO_ENABLED=%d", 0))
	env = append(env, fmt.Sprintf("GOOS=%s", "js"))
	env = append(env, fmt.Sprintf("GOARCH=%s", "wasm"))
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()

}
