# Sample Todo App Backend (Golang)

## Build and Tests
### Building the project

Project setup
```
$ go install
```

Compiles
```
$ go build -o main .
```

### Running tests

Run your unit tests
```
$ go test todo-service/service
$ go test todo-service/handler
```

## Run
```
$ docker-compose up -d
```