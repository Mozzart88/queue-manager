# Queue Manager

Simple Message Queue Manager with SQLite3 storage and a basic HTTP API.

## Table of Contents
- [Build](#build)
- [Run](#run)
- [Docker](#docker)
- [Usage](#usage)
    - [Examples](#examples)

## Build

### Build binary
```bash
go build -tags "osusergo netgo sqlite3" -ldflags "-linkmode external -extldflags -static" -o qm ./cmd/app/main.go
```

### Run tests
```bash
go test -parallel=1 -tags "integration sqlite3" ./...
```

### Create Data Base
The Data Base isn't generated automatically, so you have to create it by yourself
```bash
sqlite3 data/data.db < data/schema.sql
```

## Run
To run the service, you need to specify the `QDB_DIR` environment variable, pointing to the directory containing the database file. The file should be named `data.db`.

```bash
export QDB_DIR=$PWD/data
```

To start the Queue Manager, run:
```bash
./qm
```

## Docker
```bash
docker build -t queue-manager .
docker run -d -p 8080:8080 -v $PWD/data:/data queue-manager
```

### Docker comopse
```bash
docker compose up -d
```

## Usage
Queue Manager supports 3 message states:
- `new` - automatically assigned to all new messages
- `active` - indicates messages currently being processed
- `done` - indicates messages successfully processed (messages in this state are automatically deleted by default)

And 3 API Endpoints:
- `/publisher` - register, update, delete, or get publishers
- `/msg` - add, retrive, delete, or update the state of a single message
- `/queue` - add or retrive batch of messages

All endpoint details are described [here](./api) in the `OpenApi` format.

By default, Queue Manager listen on port `:8080`

### Examples

#### Register a publisher:
```bash
curl localhost:8080/publisher -d '\
{\
"name": "publisher name"
}'
```

Returns:
```json
{
    "msg": "{\"id\":1,\"name\":\"publisher name\"}",
    "code": 201
}
```

#### Add a new message:
```bash
curl localhost:8080/msg -d '
{
"publisher": "publisher name",
"msg": "some new message"
}'
```

Returns:
```json
{
    "msg": "{\"id\":1,\"publisher\":\"publisher name\",\"msg\":\"some new message\",\"state\":\"new\"}",
    "code": 201
}
```

#### Get a message:
```bash
curl 'localhost:8080/msg?id=1'
```

Returns:
```json
{
    "msg": "{\"id\":9,\"publisher\":\"publisher name\",\"msg\":\"some new message\",\"state\":\"new\"}",
    "code": 200
}
```

#### Update message state:
```bash
curl -X PATCH localhost:8080/msg -d '{
    "id": 1,
    "state": "active"
}'
```

Returns:
```json
{
    "msg": "ok",
    "code": 200
}
```

#### Add a batch of messages:
```bash
curl localhost:8080/queue -d '{
    "publisher": "publisher name",
    "msgs": ["content of first message", "content of second message"]
}'
```

Returns:
```json
{
    "msg": "2",
    "code": 201
}
```

#### Get a batch of messages:
```bash
curl 'localhost:8080/queue?publisher=publisher+name&state=new'
```

Returns:
```json
{
    "msg": "[{\"id\":10,\"publisher\":\"publisher name\",\"msg\":\"content of first message\",\"state\":\"new\"},{\"id\":11,\"publisher\":\"publisher name\",\"msg\":\"content of second message\",\"state\":\"new\"}]",
    "code": 200
}
```