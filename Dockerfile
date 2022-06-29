FROM golang:1.18-alpine AS BUILD

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.3

WORKDIR /app

COPY go.mod /app

COPY go.sum /app

RUN go mod download

COPY . /app

RUN swag init

RUN go build -o api

FROM alpine:3.15

COPY --from=BUILD /app/api /bin/

CMD [ "api" ] 



