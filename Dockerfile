# Set the base image to the latest official Go image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY ./go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the app
RUN go build -o svc "./cmd/packs-api"

# Run the app
CMD ["./svc"]
