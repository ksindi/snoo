NAME=snoo
VERSION=0.1.0
COMMIT=$(shell git rev-parse --short=7 HEAD)
TIMESTAMP:=$(shell date -u '+%Y-%m-%dT%I:%M:%SZ')

LDFLAGS += -X main.BuildTime=${TIMESTAMP}
LDFLAGS += -X main.BuildSHA=${COMMIT}
LDFLAGS += -X main.Version=${VERSION}

PREFIX?=${PWD}/
DOCKER=$(shell command -v docker;)
TEST_FLAGS?=-race

.PHONY: all
all: quality test

.PHONY: quality
quality:
	go vet
	go fmt
	go mod tidy
ifneq (${DOCKER},)
	docker run -v ${PWD}:/src -w /src -it golangci/golangci-lint golangci-lint run --true gocritic --true gosec --true golint --true stylecheck --exclude-use-default=false
endif

.PHONY: test
test:
	go test ${TEST_FLAGS} -coverprofile=coverage

.PHONY: clean
clean:
	rm -f ${NAME}*

.PHONY: build
build: clean build-darwin build-linux

build-%:
	GOOS=$* GOARCH=amd64 go build -ldflags '${LDFLAGS}' -o ${PREFIX}${NAME}-$*

.PHONY: docker
docker:
ifeq (${DOCKER},)
	@echo Skipping Docker build because Docker is not installed
else
	docker run --rm -i hadolint/hadolint < Dockerfile
	docker build \
	--build-arg NAME="${NAME}" \
	--build-arg VERSION="${VERSION}" \
	--build-arg COMMIT="${COMMIT}" \
	--build-arg BUILD_DATE="${TIMESTAMP}" \
	--build-arg LDFLAGS='${LDFLAGS}' \
	--tag ${NAME} .
	docker tag ${NAME} ${NAME}:${VERSION}
	docker run -it ${NAME}:${VERSION} -- -help 2>&1 | grep -F '${NAME} v${VERSION} ${TIMESTAMP} ${COMMIT}'
endif
