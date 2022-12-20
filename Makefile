server-build:
	rm -f -r build/server
	mkdir -p build/server
	cd ./backend; go build cmd/main.go; mv main ../build/server;

server-clear:
	rm -f -r build/server

server-run:
	make server-build; ./build/server/main

server-dev:
	go run backend/cmd/main.go
