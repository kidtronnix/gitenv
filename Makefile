build:
	docker run -v $$(pwd):/gitenv -w /gitenv -u $$(id -u):$$(id -g) golang:latest go build

fmt:
	docker run -v $$(pwd):/gitenv -w /gitenv -u $$(id -u):$$(id -g) golang:latest gofmt -w *.go

