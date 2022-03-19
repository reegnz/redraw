all: build

build:
	goreleaser build --snapshot --rm-dist

clean:
	rm -rf dist
