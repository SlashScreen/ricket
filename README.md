# Ricket

## What is Ricket?

Ricket is an attempt to open up the WASI platform to the Plan 9 OS, and hopefully allow more software to be used on Plan 9 with little cost to the developer.

Ricket is a frontend and toolset for [wazero](https://github.com/tetratelabs/wazero) tailored specifically for use with Plan 9.

Special thanks to Ori on the #cat-v IRC channel for making this more Plan 9-like.

Installing on a 9front system is easy:

```shell
git/clone https://github.com/SlashScreen/ricket
cd ricket
mk install
```

## Project status

Ricket is now in beta- it should work, but it may have bugs. If something goes wrong, file a bug report!

## Project goals

- [x] Be able to run WASI programs at all
- [X] Be able to run WASI programs on Plan 9
- [X] Be able to package WASI programs as standalone apps
