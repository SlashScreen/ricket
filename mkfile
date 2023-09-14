install:V Q:
    mkdir -p /$objtype/bin/ricket
	cp ricket /$objtype/bin/ricket/run
    cp package.rc /$objtype/bin/ricket/package
	cp ricket.troff /sys/man/1/ricket
    echo Ricket is now installed on your system.

clean:V Q:
    rm -rf /$objtype/bin/ricket
    rm -f /sys/man/1/ricket
    echo Bye, ricket!

build:V Q:
    echo Please note that go must be installed on the system for this to work.
    echo Attempting to build ricket from source.
    GOARCH=$objtype go build .
