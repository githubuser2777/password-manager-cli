.PHONY: build test clean

APP_NAME=passmgr

build:
	go build -o $(APP_NAME)

test:
	go test ./...

clean:
	rm -f $(APP_NAME) $(APP_NAME).exe
