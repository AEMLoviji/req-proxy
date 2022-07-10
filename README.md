# req-proxy

Simple HTTP server for proxying HTTP-requests to 3rd-party services.

## Before you run

Make sure sure you have installed `Go` at least versino 18.1.

- First timee call `go mod tidy` in order to install all required packages added to project

## How to run

By default app will run on port `8080`. If you want to choose other port then execute `go run` command with `port` argument as in example shown below, otherwise omit `-port` argument.

```sh
go run main.go -port <any port>
```

## What we have

- [x] minor `ping`<->`pong` route added.
- [x] HTTP Api with 2 routes.
  - [x] `/proxy`
  - [x] `/proxy/history`
- [x] demonstrated how concurently track requests. See [request_tracker in observer package](./observer/request_tracker.go)
- [x] `domain` & `observer` packages are covered with unit tests. See `*test.go` files in the corresponding packages
- [x] API layer covered with unit tests. See `api > proxy > api_test.go` [api_test.go](./api/proxy/api_test.go) 
- [x] graceful shutdown logic implemented. See [Start() in server.go in api package](./api/server.go)

## Request samples

cURL samples can be run from command line

### /ping

```sh
curl --request GET 'http://localhost:8080/ping'
```

### /proxy

```sh
curl --request GET 'http://localhost:8080/proxy' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "GET",
    "url": "http://yandex.com",
    "headers": {
        "Authentication": ["Basic bG9naW46cGFzc3dvcmQ="]
    }
}
'
```

### /proxy/history

```sh
curl --request GET 'http://localhost:8080/proxy/history'
```
