.PHONY: test
test:
	go test -fullpath -race -cover -bench=. ./...

lint:
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

lint-local:
	golangci-lint run

proto-gen:
	cd ./api/grpc && ./generate.sh

infra-start:
	cd ./deployments && docker compose -p boilerplate up -d

infra-stop:
	cd ./deployments && docker compose -p boilerplate stop

.PHONY: vendor
vendor:
	go mod vendor && go mod tidy
