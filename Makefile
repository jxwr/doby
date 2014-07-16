all:
	./build.sh

.PHONY: test
test:
	cd test/go && go run sdl.go -i ../sdl.d
