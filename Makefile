BIN_NAME:=mdp
GOOS=linux

build: main.go
	GOOS=$(GOOS) go build -o $(BIN_NAME)

test: main.go
	go test -v

clean:
	rm $(BIN_NAME)

.PHONY: clean
