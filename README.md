[![Go Reference](https://pkg.go.dev/badge/github.com/abu-lang/goabu.svg)](https://pkg.go.dev/github.com/bancodobrasil/featws-api)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/bancodobrasil/featws-api/blob/develop/LICENSE)

# Featws API [![About_en](https://github.com/yammadev/flag-icons/blob/master/png/BR.png?raw=true)](https://github.com/bancodobrasil/featws-api/blob/develop/README-PTBR.md)
## How to run

In order to run this project, you need to have certain prerequisites set up on your machine. These prerequisites include:

 - [Golang](https://go.dev/doc/install)
 - [Swaggo](https://github.com/swaggo/swag/blob/master/README_pt.md#come%C3%A7ando)
- [Docker](https://www.docker.com/)


To run the project, follow these steps:

- Open the terminal in the project directory and run the command `go mod tidy` to ensure that all required dependencies are installed.

- Then, run the command `swag init` to initialize Swagger and generate the necessary API documentation.

- Open Docker on your machine or a similar application to ensure that it is working properly. Then, open the terminal in the project repository and run `docker compose up`.

- Once the containers are started and running, run the command `make run` in the project terminal to start the project.

The project will run on `localhost:9007`. To access the Swagger documentation [click here](http://localhost:9007/swagger/index.html#/).

By following these steps, the project will be up and running, and you will be able to access the API documentation through Swagger.

## GoDoc

To access the GoDoc documentation, first install GoDoc on your machine. Open a terminal and type:

````
go get golang.org/x/tools/cmd/godoc
````
    
Then run the following command in the repository terminal:
    
````
godoc -http=:6060
````

GoDoc will run on `localhost:6060`. To access the GoDoc documentation, just [click here](http://localhost:6060/pkg/).

## MacOS config

The MySQL used in the Docker Compose file isn't compatible with Apple's M1 and M2 chips. To make it work, perform the following modifications:

### Arquivo .env

Create the variable FEATWS_API_MYSQL_URI and set the path as follows:

``````
FEATWS_API_MYSQL_URI=api:api@tcp(localhost:3307)/api
``````

### Docker Compose

In the Docker Compose file, it is necessary to change MySQL to MariaDB. Only the "db" section in the docker-compose file needs to be modified.

``````
  db:
    image: mariadb
    restart: always
    ports:
      - 3307:3306
    volumes:
      - ./data/mysql:/var/lib/mysql:rw
    user: mysql
    environment:
      MARIADB_ROOT_PASSWORD: "root"
      MARIADB_DATABASE: "api"
      MARIADB_USER: "api"
      MARIADB_PASSWORD: "api"
      TELEMETRY_HTTPCLIENT_TLS: "false"
      TELEMETRY_EXPORTER_JAEGER_AGENT_HOST: "localhost"
      TELEMETRY_ENVIRONMENT: local
``````