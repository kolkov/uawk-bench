.PHONY: all generate build run clean

INPUT_DIR = testdata
RESULTS_DIR = results

all: generate build run

generate:
	@echo "Generating test data..."
	@go run scripts/generate-data.go

build:
	@echo "Building benchmark tool..."
	@go build -ldflags "-s -w" -o bin/awkbench.exe ./cmd/awkbench

run:
	@echo "Running benchmarks..."
	@./bin/awkbench.exe -size 10MB -runs 5

run-small:
	@./bin/awkbench.exe -size 1MB -runs 3

run-large:
	@./bin/awkbench.exe -size 100MB -runs 5

clean:
	@rm -rf bin/*.exe testdata/* results/*

# Install AWK implementations (for local testing)
install-deps:
	@echo "Installing goawk..."
	@go install github.com/benhoyt/goawk@latest
	@echo "Done. Make sure uawk, gawk, mawk are also installed."
