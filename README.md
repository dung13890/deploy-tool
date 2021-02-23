# deploy-tool

## Demo

![](docs/images/deploy.gif?raw=true)

*Note: [init](docs/images/init.gif?raw=true) | [ping](docs/images/ping.gif?raw=true)*

## Features
- Deployment from local into remote
- Deployment on remote
- Run command for multiple remote
- Rsync multiple cluster
- Notify to chatwork, slack
- UI for deployment

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
