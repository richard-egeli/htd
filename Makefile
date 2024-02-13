build:
	go build -o build/main cmd/htd/main.go

test:
	go test cmd/htd/main.go

clean: rm -rf tmp
