
all: update test

update:
	go mod vendor

.PHONY: setup
setup:
	docker-compose up -d --force-recreate --remove-orphans

.PHONY: teardown
teardown:
	-docker-compose down --remove-orphans

test:
	 go test -cover ./...

docker-test:
	docker build -f Dockerfile.tester .
