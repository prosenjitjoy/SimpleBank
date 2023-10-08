FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
        
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY database/migration ./migration
COPY start.sh .
COPY .env .

EXPOSE 5000
CMD [ "/app/main" ]
ENTRYPOINT ["/app/start.sh"]