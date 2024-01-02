mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o heartbeats-darwin
	scp heartbeats-darwin will@192.168.1.81:

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o heartbeats-linux
	scp heartbeats-linux will@192.168.1.77:
	scp heartbeats-linux will@192.168.1.80:
	scp heartbeats-linux will@192.168.1.82: