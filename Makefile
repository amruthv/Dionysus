.PHONY: all

all:
	service alwaysbeer stop;\
	go run backend/main.go

start: 
	service alwaysbeer start

docker:
	docker build -t alwaysbeer .
