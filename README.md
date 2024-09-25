# jcd-api

## Prerequisites

* Go `1.23.1`+
* [`air`](https://github.com/air-verse/air)
  * For hot reload
* `fswatch`
  * For compile-on-change

## Getting Started

### Environment Variables
Some environment variables are required. In the case that the env var is a secret, provide a `base64` encoded string.

The project uses [`godotenv`](https://github.com/joho/godotenv) for loading environment vars during development. Provide env vars in a `.env` file at the project root:
```sh
touch .env
```

### Run the Server
```sh
make build
./bin/jcd-api
```

Hot reload:
```sh
make watch
```
