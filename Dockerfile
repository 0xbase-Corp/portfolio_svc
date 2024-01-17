# Start from the latest golang base image
FROM golang:alpine3.19 AS binary

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the migrations directory
COPY pkg/migrations /app/migrations

# Download all dependencies and tidy up
RUN go mod download && go mod tidy

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

WORKDIR /app/cmd/0xbase/

# Build the Go app
RUN go build -o ../../main-out .

FROM gcr.io/distroless/base-debian12 AS build-release-stage

COPY --from=binary /app/main-out /app/
COPY --from=binary /app/app.env ./
# Copy the migrations directory to the final stage
COPY --from=binary /app/migrations /app/migrations

EXPOSE 5050
# Command to run the executable
CMD ["/app/main-out"]