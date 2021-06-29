FROM golang:1.16-alpine AS builder

ENV PROJECT_PATH=github.com/Frosin/setserver
ENV CMD_PATH=cmd/setserver

# Set enviroment variable for Go
ENV GOPATH=/go \
    PATH="/go/bin:$PATH"

# Copy project files
WORKDIR ${GOPATH}/src/${PROJECT_PATH}
COPY . .

# Build
WORKDIR ${CMD_PATH}
RUN CGO_ENABLED=0 GOOS=linux go build -o ${GOPATH}/bin/instance .

# Init new lightweight container
FROM alpine:3.11

WORKDIR /app
ENV STORAGE_TIMEOUT=3
ENV SERVER_PORT=8080
COPY --from=builder /go/bin/instance .
CMD /app/instance
