package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Not enough arguments. Type \"ricket help\" for usage.")
		return
	}

	switch os.Args[1] {
	case "run":
		run_program()
	case "package":
		package_file()
	case "help", "?":
		help()
	default:
		fmt.Printf("Unknown command: %s", os.Args[1])
	}
}

func run_program() {
	// Check for arguments
	if len(os.Args) < 3 {
		log.Println("No path to WASM file provided.")
		return
	}

	ctx := context.Background()

	// Instantiate runtime
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Read program
	wasm, err := os.ReadFile(os.Args[2])
	if err != nil {
		log.Panicf("failed to read WASM file: %v", err)
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
		log.Panicf("failed to instantiate WASM program: %v", err)
	}
}

func package_file() {

}

func help() {
	println(`
usage:
	ricket run path [ args ... ] - run a .wasm file at <path>, passing in any arguments.
	ricket package path name bin_dir [ -o ] - package a .wasm file at <path> into a standalone program called <name> at <bin_dir>, <-o>ptionally <-o>mitting the copied ricket executable.
	ricket help | ? - open this page. Plan 9 users should instead run 'man ricket'.
	`)
}
