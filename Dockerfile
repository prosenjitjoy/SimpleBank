FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
        
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY database/migration ./database/migration
COPY .env .

CMD [ "/app/main" ]