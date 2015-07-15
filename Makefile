.PHONY: all

all:
	cd backend &&\
	go build &&\
	service alwaysbeer restart

docker:
	docker build -t alwaysbeer .
