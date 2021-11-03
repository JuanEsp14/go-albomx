	ARTIFACT_ID ?= albomx-comics
	GIN_MODE ?= debug
	LOG_LEVEL ?= debug
	GREEN=\033[0;32m
	NC=\033[0m # No Color

run-locally:
	GIN_MODE=${GIN_MODE} LOG_LEVEL=${LOG_LEVEL} EXAMPLE="" go run main.go

dependencies:
	echo "${GREEN}Pulling dependencies${NC}"
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/onsi/ginkgo/ginkgo

tests:
	echo "${GREEN}Running TESTs${NC}"
	$(HOME)/go/bin/ginkgo -r --progress --failFast  --randomizeAllSpecs --randomizeSuites --failOnPending --trace --reportFile=./junit.xml -coverpkg=./... -coverprofile=coverage.out -outputdir=./test
	go tool cover -html=./test/coverage.out -o ./test/coverage.html

build-aws-zip:
	$(MAKE) GOOS=linux GOARCH=amd64 build
	echo "${GREEN}Zipping ${ARTIFACT_ID}.zip ${NC}"
	cd bin && zip "${ARTIFACT_ID}.zip" "${ARTIFACT_ID}"

build:
	$(MAKE) swagger
	mkdir -p bin
	echo "${GREEN}Building bin/${ARTIFACT_ID} with GOOS=${GOOS} ${NC}"
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "-s -w" -o bin/${ARTIFACT_ID} main.go

swagger:
	echo "${GREEN}Bootstrapping SWAGGER ${NC}"
	$(HOME)/go/bin/swag init --parseDependency --parseInternal --parseDepth 3
