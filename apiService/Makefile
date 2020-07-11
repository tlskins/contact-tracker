.PHONY: build clean deploy gomodgen run-local

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/users users/deliveries/lambda/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/places places/deliveries/lambda/main.go
	cp ./id_rsa bin/id_rsa
	cp ./id_rsa.pub bin/id_rsa.pub

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

run-local-check-ins:
	PLACES_PORT=8007 go run cmd/server/check-ins/main.go