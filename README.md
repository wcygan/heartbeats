## Build for MacOS

```
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o heartbeats-darwin
```

## Build for Linux

```
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o heartbeats-linux
```

## Copy to remote machine 

```
scp file user@ip:
```