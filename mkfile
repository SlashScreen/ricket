install:VQ:
    mkdir -p /$objtype/bin/ricket
	cp bin/$objtype/ricket /$objtype/bin/ricket/run
    cp package.rc /$objtype/bin/ricket/package
	cp ricket.troff /sys/man/1/ricket
    echo Ricket is now installed on your system.

clean:VQ:
    rm -rf /$objtype/bin/ricket
    rm -f /sys/man/1/ricket
    echo Bye, ricket!

build:VQ:
    echo Please note that go 1.19+ must be installed on the system for this to work.
    echo Attempting to build ricket from source.
    GOARCH=$objtype go build .
