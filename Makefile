rebuild:
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
