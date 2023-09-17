# Define variables
CENTRAL_DOCKER_IMAGE = central-server
AMERICA_DOCKER_IMAGE = america-server
EUROPE_DOCKER_IMAGE = europa-server
ASIA_DOCKER_IMAGE = asia-server
OCEANIA_DOCKER_IMAGE = oceania-server
RABBITMQ_DOCKER_IMAGE = rabbitmq-server

# Define the default target
all: help

# Build the central server Docker image
docker-central:
    docker build -t $(CENTRAL_DOCKER_IMAGE) --build-arg SERVER_TYPE=central .
	docker run -d --name $(CENTRAL_DOCKER_IMAGE) -p 8081:8081 $(CENTRAL_DOCKER_IMAGE)

# Build the regional server Docker images
docker-regional:
	docker build -t $(AMERICA_DOCKER_IMAGE) --build-arg SERVER_TYPE=america ./regional/america
	docker build -t $(EUROPE_DOCKER_IMAGE) --build-arg SERVER_TYPE=europe ./regional/europe
	docker build -t $(ASIA_DOCKER_IMAGE) --build-arg SERVER_TYPE=asia ./regional/asia
	docker build -t $(OCEANIA_DOCKER_IMAGE) --build-arg SERVER_TYPE=oceania ./regional/oceania
	docker run -d --name $(AMERICA_DOCKER_IMAGE) -p 50051:50051 $(AMERICA_DOCKER_IMAGE)
	docker run -d --name $(EUROPE_DOCKER_IMAGE) -p 50052:50052 $(EUROPE_DOCKER_IMAGE)
	docker run -d --name $(ASIA_DOCKER_IMAGE) -p 50053:50053 $(ASIA_DOCKER_IMAGE)
	docker run -d --name $(OCEANIA_DOCKER_IMAGE) -p 50054:50054 $(OCEANIA_DOCKER_IMAGE)

# Build the RabbitMQ server Docker image
docker-rabbitmq:
    docker build -t $(RABBITMQ_DOCKER_IMAGE) ./rabbitmq
	docker run -d --name $(RABBITMQ_DOCKER_IMAGE) -p 5673:5673 -p 15673:15673 $(RABBITMQ_DOCKER_IMAGE)

# Usage: make help
help:
	@echo "Available targets:"
    @echo "  docker-central   - Start the Docker code for the central server"
    @echo "  docker-regional  - Start the Docker code for the regional servers"
    @echo "  docker-rabbit    - Start the Docker code for RabbitMQ"
    @echo "  help             - Display this help message"

.PHONY: all docker-central docker-regional docker-rabbit help