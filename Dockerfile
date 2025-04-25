FROM golang:1.24 AS builder

ENV QDB_DIR=/app/data

WORKDIR /app

RUN apt-get update && apt-get install -y gcc libc6-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go test -p=1 -tags "integration sqlite3" ./...

RUN go build -tags "osusergo netgo sqlite3" -ldflags "-linkmode external -extldflags -static" -o qm ./cmd/app/main.go
VOLUME [ "/app/data" ]
EXPOSE 8080

FROM alpine:latest AS prod
RUN apk add --no-cache libc6-compat

ENV QDB_DIR=/data
ENV QDB_FILE=/data/data.db

WORKDIR /app

COPY --from=builder /app/qm .

EXPOSE 8080
VOLUME [ "/data" ]

CMD [ "./qm" ]