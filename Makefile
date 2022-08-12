test:
	go test -race ./...

run-server:
	modd -f server.modd.conf