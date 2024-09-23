# jcd-api

## Prerequisites

* Go `1.23.1`+
* [`air`](https://github.com/air-verse/air)
  * For hot reload
* `fswatch`
  * For compile-on-change

## Getting Started

Regular build (w/o `air`, `fswatch`)
```sh
make build
./bin/jcd-api
```

Hot reload:
```sh
make watch
```
