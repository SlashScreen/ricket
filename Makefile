ARCHS = amd64 arm 386

build:
	@echo "Building Ricket for: $(ARCHS)"
	@$(foreach arch, $(ARCHS), mkdir -p bin/$(arch) ; GOOS=plan9 GOARCH=$(arch) go build -o bin/$(arch)/ricket ;)
	@echo "Built plan 9 executables."
