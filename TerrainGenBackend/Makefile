.PHONY: docker

all:
	go build -o TerrainGenBackend

docker:
	rm ./backend.tar
	docker build --tag terraingenbackend .
	docker save -o ./backend.tar terraingenbackend