# tutorial

## go run

```bash
go run ./cmd/server/server.go
```

## test

```bash
curl --location --request GET 'http://127.0.0.1:8888/hello'
```

## config file

## build

```
CGO_ENABLED=0 go build -trimpath -o server
```
