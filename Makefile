.PHONY: all

BUILD_FLAG =
GO := GO111MODULE=on go
GOBUILD := $(GO) build $(BUILD_FLAG)

all: test bms simulator

test: unit_test

unit_test:
	$(GO) test -v ./...

bms:
	$(GOBUILD) -o bin/bms-server ./cmd/bms.go

simulator:
	$(GOBUILD) -o bin/simulator ./cmd/simulator.go

clean:
	rm -rf bin
