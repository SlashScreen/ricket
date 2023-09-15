package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	if len(os.Args) == 0 {
		help()
	} else {
		run(os.Args[1])
	}
}

func run(path string) {
	// Instantiate runtime
	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Read program
	wasm, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("failed to read WASM file: %v\n", err)
	}

	// if present, strip off #! line
	if len(wasm) > 2 && string(wasm[0:2]) == "#!" {
		for i := 2; i < len(wasm); i++ {
			if wasm[i] == byte('\n') {
				wasm = wasm[i+1:]
				break
			}
		}
	}

	var wasmArgs []string
	if len(os.Args) > 3 {
		wasmArgs = os.Args[3:]
	}

	// Run program
	conf := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithFSConfig(wazero.NewFSConfig()).
		WithRandSource(rand.Reader).
		WithArgs(wasmArgs...)

	_, err = r.InstantiateWithConfig(ctx, wasm, conf)
	if err != nil {
		log.Panicf("failed to instantiate WASM program: %v\n", err)
	}
}

func help() {
	println(`
usage: ricket/run path [ args ... ]
run 'man ricket' for more info.
	`)
}
