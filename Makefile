.PHONY: all build clean

all: build

build:
	@./build.sh
clean:
	rm -r target