# Define variables
CENTRAL_DOCKER_IMAGE = central-server
AMERICA_DOCKER_IMAGE = america-server
EUROPA_DOCKER_IMAGE = europa-server
ASIA_DOCKER_IMAGE = asia-server
OCEANIA_DOCKER_IMAGE = oceania-server
RABBITMQ_DOCKER_IMAGE = rabbitmq-server

# Define Docker Compose file
DOCKER_COMPOSE_FILE = docker-compose.yml

# Define the default target
all: help

# Build the central server Docker image
docker-central:
    docker build -t $(CENTRAL_DOCKER_IMAGE) ./central
	docker run -d --name $(CENTRAL_DOCKER_IMAGE) $(CENTRAL_DOCKER_IMAGE)

# Build the regional server Docker image
docker-regional:
    docker build -t $(AMERICA_DOCKER_IMAGE) ./regional/america
	docker build -t $(EUROPA_DOCKER_IMAGE) ./regional/europa
	docker build -t $(ASIA_DOCKER_IMAGE) ./regional/asia
	docker build -t $(OCEANIA_DOCKER_IMAGE) ./regional/oceania
	docker run -d --name $(AMERICA_DOCKER_IMAGE) $(AMERICA_DOCKER_IMAGE)
	docker run -d --name $(EUROPA_DOCKER_IMAGE) $(EUROPA_DOCKER_IMAGE)
	docker run -d --name $(ASIA_DOCKER_IMAGE) $(ASIA_DOCKER_IMAGE)
	docker run -d --name $(OCEANIA_DOCKER_IMAGE) $(OCEANIA_DOCKER_IMAGE)

# Build the RabbitMQ server Docker image
docker-rabbitmq:
    docker build -t $(RABBITMQ_DOCKER_IMAGE) ./rabbitmq
	docker run -d --name $(RABBITMQ_DOCKER_IMAGE) $(RABBITMQ_DOCKER_IMAGE)

# Usage: make help
help:
    @echo "Available targets:"
    @echo "  docker-central   - Start the Docker container for the central server"
    @echo "  docker-regional  - Start the Docker container for the regional servers"
    @echo "  docker-rabbit    - Start the Docker container for RabbitMQ"
    @echo "  help             - Display this help message"

.PHONY: all docker-central docker-regional docker-rabbit help