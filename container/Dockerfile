# Stage 1: Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go build -o agent agent.go

## Stage 2: Final stage
#FROM alpine:latest
#
#WORKDIR /app
#
#COPY --from=builder /app/server .

ENV REDIS_MASTER_IP="10.20.1.1" \
    REDIS_MASTER_PORT="6379" \
    REDIS_MASTER_PASSWORD="" \
    REDIS_WRITE_KEY="0" \
    REDIS_LISTEN_KEY="1" \
    SQL_DB_USER="rdmaredis" \
    SQL_DB_PASSWORD="rdmaredis" \
    SQL_DB_NAME="rdmaredis" \
    SQL_DB_HOST="10.20.1.1" \
    SQL_DB_PORT="9200"

EXPOSE 8000

ENTRYPOINT ["/app/agent"]
