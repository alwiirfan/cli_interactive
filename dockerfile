# image as the base image for the build stage.
FROM golang:1.20.4-alpine AS builder

# provide information about the author, email, and a link to the GitHub repository associated with the container.
LABEL author=Mohammad_Alwi_irfani
LABEL email=alwi.irfani1927@gmail.com 
LABEL github=https://github.com/alwi09/cli_interactive

# Updating the Alpine Linux packages and installing Git.
RUN apk update && apk add --no-cache git

# Setting the current working directory inside the container
WORKDIR /home

# Copying the go.mod and go.sum files from the host directory to the current directory (/home) in the container.
COPY go.mod go.sum ./

# command to ensure that the Go dependencies stated in go.mod match the dependencies actually used by the code.
RUN go mod tidy

# Downloading the Go dependencies stated in go.mod into the current directory in the container.
RUN go mod download

# Copying all files from the current host directory to the current directory (/home) in the container.
COPY . .

# Running the go build command to compile the Go code into a binary file named cli. The compiled code comes from ./cmd/main.go.
RUN go build -o cli ./cmd/main.go

# Using the alpine:3.15 image as the base image for the next stage.
FROM alpine:3.15

# Updating the Alpine Linux packages in the next stage and installing Git.
RUN apk update && apk add --no-cache git

# Setting the current working directory inside the container to /home.
WORKDIR /home

# Copying the cli binary file generated from the previous stage to the current directory in the container. 
COPY --from=builder /home/cli .

# Copying the migrations directory inside the internal/database directory to the corresponding directory in the container.
COPY --from=builder /home/internal/database/migrations/ ./internal/database/migrations/

# Specifying the command to be executed when the container runs. In this case, the command ./cli worker will be executed when the container starts.
CMD ["./cli", "worker"]