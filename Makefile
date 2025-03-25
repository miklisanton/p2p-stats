.PHONY: all build
all: build
	docker run -d  --volume $(pwd)/p2p-stats.db:/app/p2p-stats.db stats-bot
build:
	docker build -t stats-bot .
