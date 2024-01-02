copy:
	# compile for mac
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o heartbeats-darwin

	# compile for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o heartbeats-linux

    # copy to other machines
	scp heartbeats-darwin will@192.168.1.81:
	scp heartbeats-linux will@192.168.1.77:
	scp heartbeats-linux will@192.168.1.80:
	scp heartbeats-linux will@192.168.1.82: