.PHONY: clean build deploy

clean:
	rm -rf ./bin ./vendor Gopkg.lock

build: clean
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler cmd/handler/main.go

deploy:
	sls deploy --verbose
