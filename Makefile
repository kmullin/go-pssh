BINARY=gpssh

VERSION=0.4.0

LDFLAGS=-ldflags "-X main.Version=${VERSION}"

.DEFAULT_GOAL: ${BINARY}

${BINARY}:
	CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY} ./cmd/gpssh

install:
	CGO_ENABLED=0 go install ${LDFLAGS} ./cmd/gpssh

clean:
	@rm -f ${BINARY}

.PHONY: clean install
