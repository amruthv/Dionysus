.PHONY: all

all:
	go run backend/main.go

docker:
	docker build -t alwaysbeer .
