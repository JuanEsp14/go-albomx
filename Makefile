	ARTIFACT_ID ?= albomx-comics
	GREEN=\033[0;32m
	NC=\033[0m # No Color

run-locally:
 	DATABASE_URL="" PRIVATE_KEY="" API_KEY="" go run main.go

dependencies:
	echo "${GREEN}Pulling dependencies${NC}"
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/onsi/ginkgo/ginkgo

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
