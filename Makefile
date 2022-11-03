.PHONY: *

build:
	docker-compose build

up:
	docker-compose up -d --force-recreate --remove-orphans

down:
	docker-compose down

bash-%:
	docker-compose exec $* bash

logs-%:
	docker-compose logs -f $*


