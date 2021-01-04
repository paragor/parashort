DOCKER_COMPOSE_BASE=deployments/docker_compose
DOCKER_COMPOSE_ARGS=-f ${DOCKER_COMPOSE_BASE}/docker-compose.yml --env-file ${DOCKER_COMPOSE_BASE}/.env

.PHONY: up
up:
	@echo "--> Up..."
	docker-compose ${DOCKER_COMPOSE_ARGS} up

.PHONY: run
run:
	@echo "--> Run..."
	docker-compose ${DOCKER_COMPOSE_ARGS} up -d

.PHONY: stop
stop:
	@echo "--> Stop..."
	docker-compose ${DOCKER_COMPOSE_ARGS} down

.PHONY: build
build:
	@echo "--> Building Docker Image..."
	docker-compose ${DOCKER_COMPOSE_ARGS} build

.PHONY: dockerhub
dockerhub:
	@echo "--> Docker hub..."
	docker build -f build/Dockerfile.backend -t paragor/parashort:latest . && docker push paragor/parashort:latest
