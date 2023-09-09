package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const arch = "amd64"

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
		log.Panicf("failed to read WASM file: %v\n", err)
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

func package_file() {
	// Check for arguments
	if len(os.Args) < 5 {
		log.Println("Improper arguments. See `ricket help` or `man ricket`.")
		return
	}

	wasm_path := os.Args[2]
	bin_dir := os.Args[4]
	program_name := os.Args[3]

	err := os.MkdirAll(bin_dir+"/"+program_name, fs.ModeAppend)
	if err != nil {
		fmt.Printf("Error while making destination directory: %s\n", err)
		return
	}

	{ // Step 1: Copy wasm file
		_, wasm_filename := path.Split(wasm_path) // make sure we only get the wasm bit
		dst := fmt.Sprintf("%s/%s/%s", bin_dir, program_name, wasm_filename)
		dest_file, err := os.Create(dst)
		if err != nil {
			fmt.Printf("Error while copying wasm file: %s\n", err)
			return
		}
		wasm_file, err := os.Open(wasm_path)
		if err != nil {
			fmt.Printf("Error while copying wasm file: %s\n", err)
			return
		}

		io.Copy(dest_file, wasm_file)
	}

	{ // Step 2: Copy ricket file if necessary
		omit := len(os.Args) == 6 && os.Args[5] == "-o"
		if !omit {
			ricket_path, err := os.Executable()
			if err != nil {
				fmt.Printf("Error while copying ricket file: %s\n", err)
				return
			}
			ricket_exec, err := os.Open(ricket_path)
			if err != nil {
				fmt.Printf("Error while copying ricket file: %s\n", err)
				return
			}
			dest_file, err := os.Create(fmt.Sprintf("%s/%s/ricket", bin_dir, program_name))
			if err != nil {
				fmt.Printf("Error while copying ricket file: %s\n", err)
				return
			}

			io.Copy(dest_file, ricket_exec)
		}
	}

	{ // Step 3: Write RC file
		dst := fmt.Sprintf("%s/%s/%s", bin_dir, program_name, program_name)
		rc, err := os.Create(dst)
		if err != nil {
			fmt.Printf("Error while creating rc file: %s\n", err)
			return
		}
		rc.Write([]byte(format_rc(wasm_path)))
	}

	{ // Step 4: Write install file
		output := format_install(program_name)
		dst, err := os.Create(fmt.Sprintf("%s/mkfile", bin_dir))
		if err != nil {
			fmt.Printf("Error while writing install file: %s\n", err)
			return
		}
		dst.Write([]byte(output))
	}
}

func help() {
	println(`
usage:
	ricket run path [ args ... ] - run a .wasm file at <path>, passing in any arguments.
	ricket package path name bin_dir [ -o ] - package a .wasm file at <path> into a standalone program called <name> at <bin_dir>, <-o>ptionally <-o>mitting the copied ricket executable.
	ricket help | ? - open this page. Plan 9 users should instead run 'man ricket'.
	`)
}

func format_rc(full_path string) string {
	_, path := path.Split(full_path)
	return fmt.Sprintf(`#!/bin/rc
ricket run %s $*
`, path)
}

func format_install(name string) string {
	return fmt.Sprintf(`
ARCH = %s
USER = glenda

install:V:
	cp %s /$ARCH/%s/bin
	echo "bind -b /$ARCH/%s/bin /bin" >> /usr/$USER/lib/profile
	`, arch, name, name, name)
} // TODO: Other architectures
