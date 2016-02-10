.PHONY: clean stop rm docker_build build app_build dev test shell up

# Important:
PROJECT_NAME=riakapi

# Do not touch
DC_BIN=docker-compose
DOCKER_COMPOSE_CMD_DEV=${DC_BIN} -p ${PROJECT_NAME} -f ../docker-compose.yml -f ./docker-compose.dev.yml


TEST_PACKAGES=./config

default:build

# Removes all the images
clean:
	rm -rf ./bin
	docker images -q --filter "dangling=true"|xargs docker rmi -f

# Stops all the cointainers
stop:
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} stop

# Removes all the containers
rm: stop
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} rm -f

# Builds docker images
docker_build:
	  cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} build

# Builds all the ecosystem
build: docker_build
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} run --rm app /bin/bash -ci  "./environment/dev/build.sh"; \
		${DOCKER_COMPOSE_CMD_DEV} stop; \
		${DOCKER_COMPOSE_CMD_DEV} rm -f

# Builds the application binary
app_build: docker_build
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} run --rm app /bin/bash -ci  "./environment/dev/build.sh; \
			go build -o bin/${PROJECT_NAME}d ./cmd/riakapi/main.go"; \
		${DOCKER_COMPOSE_CMD_DEV} stop; \
		${DOCKER_COMPOSE_CMD_DEV} rm -f

# Runs the applications
dev: docker_build
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} run --rm --service-ports app /bin/bash -ci "./environment/dev/build.sh;go run ./cmd/riakapi/main.go"; \
		${DOCKER_COMPOSE_CMD_DEV} stop; \
		${DOCKER_COMPOSE_CMD_DEV} rm -f

# Runs test suite
test: docker_build
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} run --rm app /bin/bash -ci "./environment/dev/build.sh;go test ${TEST_PACKAGES} -v"; \
		${DOCKER_COMPOSE_CMD_DEV} rm -f

# Loads a shell without binded ports
shell: docker_build
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} run --rm app /bin/bash

# Loads a shell with binded ports
up: docker_build
	cd environment/dev && \
		${DOCKER_COMPOSE_CMD_DEV} run --rm --service-ports app /bin/bash -ci "./environment/dev/build.sh;/bin/bash";

authors:
	-git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf > ./AUTHORS
