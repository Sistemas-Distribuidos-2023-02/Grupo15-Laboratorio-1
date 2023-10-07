# Define variables
CENTRAL_DOCKER_IMAGE = central-server
AMERICA_DOCKER_IMAGE = america-server
EUROPE_DOCKER_IMAGE = europa-server
ASIA_DOCKER_IMAGE = asia-server
OCEANIA_DOCKER_IMAGE = oceania-server
RABBITMQ_DOCKER_IMAGE = rabbitmq-server

AMERICA_PORT= 50051
ASIA_PORT= 50052
EUROPE_PORT= 50053
OCEANIA_PORT= 50054

# Define the default target
all: help

# Build the central server Docker image
docker-central:
	docker build -t $(CENTRAL_DOCKER_IMAGE) --build-arg SERVER_TYPE=central .
	docker run -d --name $(CENTRAL_DOCKER_IMAGE) -p 8081:8081 $(CENTRAL_DOCKER_IMAGE)

# Build the regional server Docker images
docker-regional:
	@case "$(SERVER_TYPE)" in \
		america) PORT=$(AMERICA_PORT);; \
		asia) PORT=$(ASIA_PORT);; \
		europa) PORT=$(EUROPE_PORT);; \
		oceania) PORT=$(OCEANIA_PORT);; \
		*) echo "Invalid SERVER_TYPE argument. Use 'america', 'asia', 'europa', or 'oceania'."; exit 1;; \
	esac; \
	docker build -t $(SERVER_TYPE)-server --build-arg SERVER_TYPE=$(SERVER_TYPE) .
	docker run -d --name $(SERVER_TYPE)-server -p $$PORT:$$PORT $(SERVER_TYPE)-server

# Build the RabbitMQ server Docker image
docker-rabbitmq:
	docker build -t $(RABBITMQ_DOCKER_IMAGE) -f rabbitmq/dockerfile ./rabbitmq
	docker run -d --name $(RABBITMQ_DOCKER_IMAGE) -p 5673:5673 -p 15673:15673 $(RABBITMQ_DOCKER_IMAGE)

# Usage: make help
help:
	@echo "Available targets:"
	@echo "  docker-central   - Iniciar el codigo Docker para el servidor central"
	@echo "  docker-regional SERVER_TYPE={america,asia,europa,oceania}  - Iniciar el codigo Docker para el servidor regional especificado"
	@echo "  docker-rabbit    - Iniciar el codigo Docker para el servidor RabbitMQ"
	@echo "  help             - Pide ayuda"

.PHONY: all docker-central docker-regional docker-rabbit help
