build:
	go build -o bin/fs

run: build
	k./bin/fs

test:
	go test ./... -vk