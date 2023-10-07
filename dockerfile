# Use the specified base image
ARG BASE_IMAGE=golang:1.18
FROM ${BASE_IMAGE} AS builder

ARG CENTRAL_PORT=8081
ARG AMERICA_PORT=50051
ARG ASIA_PORT=50052
ARG EUROPE_PORT=50053
ARG OCEANIA_PORT=50054

ARG SERVER_TYPE

# Set the working directory inside the container
WORKDIR /app

# Copy the parent directory's go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download and install Go dependencies
RUN go mod download

# Copy the rest of your application code to the container
COPY . .

CMD if [ "$SERVER_TYPE" = "central" ]; then \
        PORT=$CENTRAL_PORT; \
        cd /app/central; \
        go build -o central-server; \
        ./central-server; \
    elif [ "$SERVER_TYPE" = "america" ]; then \
        PORT=$AMERICA_PORT; \
        cd /app/regional/america; \
        go build -o america-server; \
        ./america-server; \
    elif [ "$SERVER_TYPE" = "asia" ]; then \
        PORT=$ASIA_PORT; \
        cd /app/regional/asia; \
        go build -o asia-server; \
        ./asia-server; \
    elif [ "$SERVER_TYPE" = "europa" ]; then \
        PORT=$EUROPE_PORT; \
        cd /app/regional/europa; \
        go build -o europa-server; \
        ./europa-server; \
    elif [ "$SERVER_TYPE" = "oceania" ]; then \
        PORT=$OCEANIA_PORT; \
        cd /app/regional/oceania; \
        go build -o oceania-server; \
        ./oceania-server; \
    else \
        echo "Invalid SERVER_TYPE argument. Use 'central', 'america', 'asia', 'europe', or 'oceania'."; \
    fi
