#!/bin/bash

docker-compose up -d adminer

go build -o api && ./api