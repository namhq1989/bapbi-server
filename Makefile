#!bin/bash

install-tools:
	@echo installing tools
	@go install	google.golang.org/protobuf/cmd/protoc-gen-go
	@go install	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@echo done

generate:
	@echo running code generation
	go generate ./...
	@echo done

run:
	doppler run -- go run cmd/*.go

mock-gen:
	mockgen -source=internal/genproto/authpb/api_grpc.pb.go -destination=internal/mock/grpc/auth_client.go -package=mockgrpc
	mockgen -source=internal/genproto/userpb/api_grpc.pb.go -destination=internal/mock/grpc/user_client.go -package=mockgrpc

test:
	doppler run -c test -- gotestsum --format testname ./pkg/...

test-coverage:
	doppler run -c test -- gotestsum --format testname -- -coverprofile=coverrage.out ./pkg/... && \
 	go tool cover -html=coverrage.out -o coverage.html

test-debug:
	doppler run -c test -- go test -v ./pkg/...
