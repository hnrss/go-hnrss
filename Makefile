src := $(wildcard *.go)

all: go-hnrss_linux_amd64 upload

go-hnrss_linux_amd64: $(src)
	gox -osarch=linux/amd64

upload:
	scp go-hnrss_linux_amd64 hnrss@hnrss.org:~

clean:
	rm -f go-hnrss_linux_amd64

.PHONY: upload clean
