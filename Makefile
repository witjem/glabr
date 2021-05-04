B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

build: info
	- cd cmd/grabr && go build -ldflags "-X main.revision=$(REV) -s -w" -o ../../dist/grabr

info:
	- @echo "revision $(REV)"