# Stage 1: Build the Go application
FROM golang:1.23.0 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum files are unchanged
RUN go mod download

# Copy the rest of the application code
COPY . ./

# Copy environment file
COPY .env.example ./.env

# install pdf
# RUN apt install wkhtmltopdf -y

# Build the Go application
RUN make build

# Stage 2: Create a lightweight image to run the application
FROM golang:1.23.0

# Set working directory for the application
WORKDIR /app


RUN apt-get update && \
    apt-get install  wkhtmltopdf -y
    
# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/production ./bin/production
COPY --from=builder /app/.env ./.env

# Expose the port the application runs on
EXPOSE 8600

# Run the application
ENTRYPOINT ["./bin/production"]
