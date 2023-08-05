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
	// Check for arguments
	if len(os.Args) < 2 {
		log.Panicf("No path to WASM file given. Usage: ricket [path]")
	}

	ctx := context.Background()

	// Instantiate runtime
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Read program
	wasm, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Panicf("failed to read WASM file: %v", err)
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
		WithRandSource(rand.Reader)

	_, err = r.InstantiateWithConfig(ctx, wasm, conf)
	if err != nil {
		log.Panicf("failed to instantiate WASM program: %v", err)
	}
}
