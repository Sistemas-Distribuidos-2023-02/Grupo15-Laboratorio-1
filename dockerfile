ARG BASE_IMAGE=golang:1.18
FROM ${BASE_IMAGE} AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the parent directory's go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download and install Go dependencies
RUN go mod download

# Copy the rest of your application code to the container
COPY . .

# Expose the port for the container (default to 8080 if not provided)
ARG PORT=8080
EXPOSE ${PORT}

# Set the default ports for the central and regional servers
ARG CENTRAL_PORT=8081
ARG AMERICA_PORT=50051
ARG ASIA_PORT=50052
ARG EUROPE_PORT=50053
ARG OCEANIA_PORT=50054

# Command to run the appropriate server or RabbitMQ
ARG SERVER_TYPE
CMD if [ "$SERVER_TYPE" = "central" ]; then \
        PORT=$CENTRAL_PORT; \
        cp -r central/ .; \
        go build -o central-server; \
        RUN go clean; \ 
        ./central-server; \
    elif [ "$SERVER_TYPE" = "america" ]; then \
        PORT=$AMERICA_PORT; \
        cp -r regional/america/ .; \
        go build -o america-server; \
        RUN go clean; \ 
        ./america-server; \
    elif [ "$SERVER_TYPE" = "asia" ]; then \
        PORT=$ASIA_PORT; \
        cp -r regional/asia/ .; \
        go build -o asia-server; \
        RUN go clean; \ 
        ./asia-server; \
    elif [ "$SERVER_TYPE" = "europe" ]; then \
        PORT=$EUROPE_PORT; \
        cp -r regional/europe/ .; \
        go build -o europe-server; \
        RUN go clean; \ 
        ./europe-server; \
    elif [ "$SERVER_TYPE" = "oceania" ]; then \
        PORT=$OCEANIA_PORT; \
        cp -r regional/oceania/ .; \
        go build -o oceania-server; \
        RUN go clean; \ 
        ./oceania-server; \
    else \
        echo "Invalid SERVER_TYPE argument. Use 'central', 'regional', or 'rabbitmq'."; \
    fi

# Final image for the server
FROM ${BASE_IMAGE}

# Copy the built server or RabbitMQ from the builder stage
COPY --from=builder /app/ /app/

# Expose the port specified in the CMD instruction
EXPOSE ${PORT}