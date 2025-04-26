build-cmd:
	go build -o qwe .
	sudo mv qwe /usr/local/bin

cmd-first-build:
	go mod download
	go build .
