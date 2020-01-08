all: bin
	go build -v -o ./bin ./cmd/...

bin:
	mkdir bin

clean:
	rm -v ./bin/*

install:
	go install -v ./cmd/...

test:
	go test -v -cover ./...
