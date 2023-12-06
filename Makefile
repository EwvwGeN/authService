genproto:
	protoc -I proto/authService proto/authService/auth.proto \
	--go_out=./proto/gen/go --go_opt=paths=source_relative \
	--go-grpc_out=./proto/gen/go/ --go-grpc_opt=paths=source_relative