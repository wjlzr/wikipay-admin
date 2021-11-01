PROJECT:=wikipay-admin

.PHONY: build
build:
	CGO_ENABLED=0 go build -o wikipay-admin main.go
build-sqlite:
	go build -tags sqlite3 -o wikipay-admin main.go
#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t wikipay-admin:latest
