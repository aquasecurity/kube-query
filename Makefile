clean:
	rm -f kube-query

build:
	go build -o kube-query kube-query.go

all:
	make clean
	make build