run:
	GIN_MODE=release ./main 2>> verbose.log

dev:
	go run main.go

build:
	go build main.go

clean:
	rm data.db
	rm verbose.log