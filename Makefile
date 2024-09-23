
GO_SRC_DIR = gosrc
BIN_DIR = bin
GO_BIN_NAME = jcd-api
GO_BIN_PATH = ${BIN_DIR}/${GO_BIN_NAME}

build:
	go build -o $(GO_BIN_PATH) ${GO_SRC_DIR}/main.go
run:
	./$(GO_BIN_PATH)
watch: build
	air --build.cmd "make build" --build.bin "make run"
	# fswatch -r ./${GO_SRC_DIR} | xargs -n1 -I{} make build
watch-build: build
	fswatch -r ./${GO_SRC_DIR} | xargs -n1 -I{} make build