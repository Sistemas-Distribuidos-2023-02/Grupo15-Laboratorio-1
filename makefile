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
	@echo "SERVER_TYPE is set to: $(SERVER_TYPE)"
	@if [ "$(SERVER_TYPE)" = "america" ]; then \
		docker build -t $(AMERICA_DOCKER_IMAGE) --build-arg SERVER_TYPE=america .; \
		docker run -d --name $(AMERICA_DOCKER_IMAGE) -p $(AMERICA_PORT):$(AMERICA_PORT) $(AMERICA_DOCKER_IMAGE); \
	elif [ "$(SERVER_TYPE)" = "asia" ]; then \
		docker build -t $(ASIA_DOCKER_IMAGE) --build-arg SERVER_TYPE=asia .; \
		docker run -d --name $(ASIA_DOCKER_IMAGE) -p $(ASIA_PORT):$(ASIA_PORT) $(ASIA_DOCKER_IMAGE); \
	elif [ "$(SERVER_TYPE)" = "europa" ]; then \
		docker build -t $(EUROPE_DOCKER_IMAGE) --build-arg SERVER_TYPE=europa .; \
		docker run -d --name $(EUROPE_DOCKER_IMAGE) -p $(EUROPE_PORT):$(EUROPE_PORT) $(EUROPE_DOCKER_IMAGE); \
	elif [ "$(SERVER_TYPE)" = "oceania" ]; then \
		docker build -t $(OCEANIA_DOCKER_IMAGE) --build-arg SERVER_TYPE=oceania .; \
		docker run -d --name $(OCEANIA_DOCKER_IMAGE) -p $(OCEANIA_PORT):$(OCEANIA_PORT) $(OCEANIA_DOCKER_IMAGE); \
	else \
		echo "Invalid SERVER_TYPE argument. Use 'america', 'asia', 'europa', or 'oceania'."; \
		exit 1; \
	fi


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
