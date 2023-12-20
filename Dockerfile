# Start from the latest golang base image
FROM golang:alpine3.19 AS binary

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .


WORKDIR /app/cmd/0xbase/

# Build the Go app
RUN go build -o ../../main-out .


FROM gcr.io/distroless/base-debian12 AS build-release-stage

COPY --from=binary /app/main-out /app/
COPY --from=binary /app/app.env ./

EXPOSE 5050
# Command to run the executable
CMD ["/app/main-out"]
