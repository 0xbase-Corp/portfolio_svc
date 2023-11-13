# portfolio_svc

This repositry is a written in Golang 

The Portfolio Backend service is responsible for handling all the requests made by **portfolio_ui** service. 

## How to Run
1. Copy all the contents from `example.env` into a new file `staging.env`  and replace all *XXXX* with the correct values. 
> Donot commit **staging.env** file
2. Run `go run cmd/main/main.go` from the root directory

Navigatge to `localhost:5050/healthy` for health check of the server

## Directory Structure

Repository Layout is based on golang community recomneded best practices. More on it [here](https://github.com/golang-standards/project-layout) 

## Libraries 

- gin is a highly scalable, light weight http server. 
- gorm to connect to a relational database 

## Adding a new ENV variable
1. add it to app.env
1. add the varaible to struct in `configs/env-config.go`
