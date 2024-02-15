BINARY=gpssh

VERSION=0.2.4

LDFLAGS=-ldflags "-X main.Version=${VERSION}"

.DEFAULT_GOAL: ${BINARY}

${BINARY}:
	go build ${LDFLAGS} -o ${BINARY} ./cmd/gpssh

install:
	go install ${LDFLAGS} ./...

clean:
	@rm -f ${BINARY}

.PHONY: clean install