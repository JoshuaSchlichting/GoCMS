
BINARY_NAME=gocms
MAIN_PATH=.
BUILD_PATH=./bin
COVERAGE_DIR=./out

all: test build

build:
	go build -o $(BUILD_PATH)/$(BINARY_NAME) -v $(MAIN_PATH)

test:
	@if [ ! -d $(COVERAGE_DIR) ]; then \
        mkdir -p $(COVERAGE_DIR); \
    fi
	go test -coverprofile=out/coverage.out -v ./...
	go tool cover -html=out/coverage.out -o out/coverage.html

clean:
	go clean
	rm -rf ./bin
	rm -rf ./out

run:
	go build -o $(BUILD_PATH)/$(BINARY_NAME) -v $(MAIN_PATH)
	$(BUILD_PATH)/$(BINARY_NAME)
