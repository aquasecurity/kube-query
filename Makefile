clean:
	rm -f kquery

build:
	go build -o kquery kube-query.go

all:
	make clean
	make build