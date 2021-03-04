# deploy-tool

# Install deploy tool

```sh
curl -sf https://gobinaries.com/dung13890/deploy-tool | PREFIX=/tmp sh
sudo mv /tmp/deploy-tool /usr/local/bin/doo
```

install doo with version

```sh
curl -sf https://gobinaries.com/dung13890/deploy-tool@1.0.1 | PREFIX=/tmp sh
sudo mv /tmp/deploy-tool /usr/local/bin/doo
```

Run doo

```sh
doo init
```

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

# For Developer

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
# Init file config.yml
go run main.go init
```
