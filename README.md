# portfolio_svc

This repositry is a written in Golang

The Portfolio Backend service is responsible for handling all the requests made by **portfolio_ui** service.

## How to Run

1. Copy all the contents from `example.env` into a new file `app.env` and replace all _XXXX_ with the correct values.
   > Donot commit **app.env** file
2. After installing dependencies run `go run cmd/0xbase/main.go` from the root directory

### Alternatively

The project can also be started by simply

```
    docker-compose up
```

Navigatge to `localhost:5050/healthy` for health check of the server

## Directory Structure

Repository Layout is based on golang community recomneded best practices. More on it [here](https://github.com/golang-standards/project-layout)

## Libraries

- gin is a highly scalable, light weight http server.
- gorm to connect to a relational database

## Adding a new ENV variable

1. add it to example.env
1. add the varaible to struct in `configs/env-config.go`

## To Update Swagger

```
swag init -g cmd/0xbase/main.go -o docs/
```

Once the server is up, the swagger UI will be available at http://localhost:5050/swagger/index.html

## Air - Live reload for Go apps

check for install (air)[https://github.com/cosmtrek/air]

replace these command to .ait.toml file with and run `air`

```
  bin = "./tmp/0xbase"
  cmd = "go build -o ./tmp/0xbase ./cmd/0xbase"
```

## Migrations

```sh
    make migrateup
    make migratedown
```
