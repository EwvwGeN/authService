genproto:
	protoc -I proto/authService proto/authService/auth.proto \
	--go_out=./proto/gen/go --go_opt=paths=source_relative \
	--go-grpc_out=./proto/gen/go/ --go-grpc_opt=paths=source_relative
build:
	CGO_ENABLED=0 GOOS=linux go build -o serverMain ./cmd/server/