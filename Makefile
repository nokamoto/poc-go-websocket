
fmt:
	gofmt -d .
	gofmt -w .

go:
	dep ensure
	go test .

docker: fmt go
	docker-compose down
	docker-compose build
