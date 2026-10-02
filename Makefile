build-cmd:
	go build -o vin .
	sudo mv vin /usr/local/bin

cmd-first-build:
	go mod download
	go build .
