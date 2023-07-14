BINARY=gpssh

VERSION=0.0.1
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

.DEFAULT_GOAL: ${BINARY}

${BINARY}:
	go build ${LDFLAGS} -o ${BINARY} ./...

install:
	go install ${LDFLAGS} ./...

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

.PHONY: clean install
