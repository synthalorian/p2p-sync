.PHONY: build tidy clean

build:
	go build -o p2p-sync ./cmd/p2p-sync

tidy:
	go mod tidy

clean:
	rm -f p2p-sync
