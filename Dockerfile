#BUILD
FROM golang:1.21.0-alpine3.17 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o main main.go

RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz


#RUN 
FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/migrate.linux-amd64 ./migrate

COPY config.yml .

COPY db/migration ./migration

COPY start.sh .

RUN chmod a+x ./start.sh 

EXPOSE 8080

CMD ["/app/main"]

ENTRYPOINT [ "./start.sh" ]