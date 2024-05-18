protoc:
	protoc --go_out=. --go-grpc_out=. ./pkg/infrastructure/pb/*.proto
server:
	go run cmd/main.go
gofmt:
	gofmt -w .
build:
	go build -o ./cmd/ciaoChatNrelExec ./cmd/main.go