clean:
	rm -f kube-query

build:
	go build -o kube-query kube-query.go

all:
	make clean
	make build

tests:
	GO111MODULE=on go test -v -short -race -timeout 30s -coverprofile=coverage.txt -covermode=atomic ./...