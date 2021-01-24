# deploy-tool

## Run with docker for development
```sh
docker-compose up -d
docker-compose exec tool bash
```

## Setup project for development
```sh
# Inside docker
go get
go mod vendor
```

## Setup and running for development

```sh
# file config.yml
go run main.go
```
