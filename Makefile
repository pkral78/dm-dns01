# Static, dependency-free builds for container use (no libc, runs in scratch/alpine).
LDFLAGS := -s -w
BUILD := CGO_ENABLED=0 go build -a -ldflags="$(LDFLAGS)"

.PHONY: all arm64 amd64 test clean
all: arm64 amd64
arm64:
	GOOS=linux GOARCH=arm64 $(BUILD) -o dm-dns01-arm64 .
amd64:
	GOOS=linux GOARCH=amd64 $(BUILD) -o dm-dns01-amd64 .
test:
	go test ./...
clean:
	rm -f dm-dns01-arm64 dm-dns01-amd64
