.PHONY: all
all:
	@mkdir -p build
	go build -v -o ./build ./cmd/...

.PHONY: clean
clean:
	rm -v ./build/*

.PHONY: install
install:
	go install -v ./cmd/...

.PHONY: test
test:
	go test -v -cover ./...
