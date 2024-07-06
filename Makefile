.PHONY: all build clean install

all: build

build:
	@./build.sh
clean:
	rm -r target
install:
	scp target/rest-arm64-linux nick@10.21.0.190:/home/nick/