.PHONY: all

all:
	cd backend &&\
	go build &&\
	cd .. &&\
	cp upstart/alwaysbeer.conf /etc/init/ &&\
	service alwaysbeer restart

docker:
	docker build -t alwaysbeer .
