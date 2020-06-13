.PHONY: build clean deploy gomodgen run-local

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/users users/transport/lambda/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/places places/transport/lambda/main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

run-local-users:
	USERS_PORT=8005 go run cmd/server/users/main.go

run-local-places:
	PLACES_PORT=8006 go run cmd/server/places/main.go