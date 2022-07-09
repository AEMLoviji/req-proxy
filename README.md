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

- [x] HTTP Api with 2 routes.
  - [x] `/proxy`
  - [x] `/proxy/logs`
- [x] demonstrated how concurently track requests. See [request_tracker in observer package](./observer/request_tracker.go) in tracker package
- [x] `domain` & `observer` packages are covered with tests. See `*test.go` files in the corresponding packages

## Request samples

cURL samples can be run from command line

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

### /proxy/logs

```sh
curl --request GET 'http://localhost:8080/proxy/logs'
```
