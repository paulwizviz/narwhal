ARG UBUNTU_VER=3.18
ARG GO_VER=1.21-bullseye

# Stage 2: Plugin build stage (Debian for CGO compatibility - if needed)
FROM golang:${GO_VER} AS builder

ARG PROTOBUF_VERSION=25.1

WORKDIR /opt

# Set necessary environment variables for Go modules
# ENV GO111MODULE=on
ENV GOPATH=/go

RUN apt-get update && apt-get install -y --no-install-recommends \
    unzip \
    wget 

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip && \
    unzip protoc-25.1-linux-x86_64.zip && \
    cp /opt/bin/protoc /usr/local/bin/protoc && \
    cp -r /opt/include /usr/local

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Stage 3: Final image (Debian base with protoc and plugins)
FROM debian:bullseye-slim

# Copy protoc binary
COPY --from=builder /usr/local/bin/protoc /usr/local/bin/protoc
COPY --from=builder /usr/local/include /usr/local/include

# Copy the plugin binaries
COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/protoc-gen-go
COPY --from=builder /go/bin/protoc-gen-go-grpc /usr/local/bin/protoc-gen-go-grpc

# Install necessary libraries for protoc to run
RUN apt-get update && apt-get install -y --no-install-recommends libprotobuf-dev

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/protoc"]

# Add metadata
LABEL description="Protoc with Go plugins (CGO compatible)"