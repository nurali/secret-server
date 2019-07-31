.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o .out/secret-service

.PHONY: build-ci
build-ci: build
	docker build -f Dockerfile -t secret-service:latest .

.PHONY: test
test:
	go test ./... -v -count=1

.PHONY: clean
clean:
	-rm -rf .out/
