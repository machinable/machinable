install:
	docker run --rm -it -v $(shell pwd):/go/src/bitbucket.org/nsjostrom/machinable -w /go/src/bitbucket.org/nsjostrom/machinable \
    instrumentisto/glide install

rebuild: install
	docker-compose build --no-cache

build:
	docker-compose build

up:
	docker-compose up -d

stop:
	docker-compose stop

remove:
	docker-compose rm -f

clean:
	docker-compose down --rmi all -v --remove-orphans
