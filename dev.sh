#!/bin/bash

docker-compose up -d mongo-express

go build -o api && ./api