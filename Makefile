install:
	docker run --rm -it -v $(shell pwd):/go/src/github.com/anothrnick/machinable -w /go/src/github.com/anothrnick/machinable \
    instrumentisto/glide install

rebuild: install
	docker-compose build --no-cache

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

delvol:
	docker volume rm machinable_db-data

stop:
	docker-compose stop

remove:
	docker-compose rm -f

clean:
	docker-compose down --rmi all -v --remove-orphans
