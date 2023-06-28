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

//go:embed test.wasm
var guestWasm []byte

func main() {
	ctx := context.Background()

	// Instantiate runtime
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

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

	_, err := r.InstantiateWithConfig(ctx, guestWasm, conf)
	if err != nil {
		log.Panicf("failed to instantiate WASM program: %v", err)
	}
}
