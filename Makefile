.PHONY: build test vet run clean

BINARY=barcodekit

build:
	CGO_ENABLED=0 go build -o $(BINARY) ./cmd/barcodekit

test:
	go test -v ./...

vet:
	go vet ./...

run: build
	./$(BINARY) -no-auth

clean:
	rm -f $(BINARY)
