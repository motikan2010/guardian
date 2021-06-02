GOFILES = $(shell find . -name '*.go')
RULEFILES = $(shell find ./crs -name '*.conf')
RULEDATAFILES = $(shell find ./crs -name '*.data')
JSONFILES = = $(shell find . -name '*.json')

default: build

build: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o docker/guardian/workdir/guardian .
	cp appsettings.live.json docker/guardian/workdir/
	docker-compose build

build-native: $(GOFILES)
	go build -o workdir/native-guardian .

up:
	docker-compose up

up-d:
	docker-compose up -d

exec-db:
	docker-compose exec db bash

exec-waf:
	docker-compose exec guardian bash

stop:
	docker-compose stop

down:
	docker-compose down

ps:
	docker-compose ps

logs:
	docker-compose logs -f guardian
