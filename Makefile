BINARY=gpssh

VERSION=0.3.0

LDFLAGS=-ldflags "-X main.Version=${VERSION}"

.DEFAULT_GOAL: ${BINARY}

${BINARY}:
	go build ${LDFLAGS} -o ${BINARY} ./cmd/gpssh

install:
	go install ${LDFLAGS} ./cmd/gpssh

clean:
	@rm -f ${BINARY}

.PHONY: clean install
