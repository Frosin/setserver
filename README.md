generation:
    protoc -I=proto --go_out=internal/api --go-grpc_out=internal/api proto/api.proto
build:
    docker build -t setserver .
    docker container run --name setserver -d -p 8080:8080 setserver
evans testing:
    evans ./proto/api.proto -p 8080
