generation:
    protoc -I=proto --go_out=internal/api --go-grpc_out=internal/api proto/api.proto