build-cmd:
	go build -o svc .
	sudo mv svc /usr/local/bin

cmd-first-build:
	go mod download
	go build .
