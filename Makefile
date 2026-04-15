.PHONY: build test clean

build:
	go build -o gitignore-merge ./cmd/gitignore-merge/

test:
	go test ./...

clean:
	rm -f gitignore-merge
