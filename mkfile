ARCH = amd64

install:V:
	cp ricket /$ARCH/bin/ricket
	cp ricket.troff /sys/man/1/ricket

clean:V:
    rm /$ARCH/bin/ricket
    rm /sys/man/1/ricket
