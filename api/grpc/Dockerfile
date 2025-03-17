ARG GOLANG_VERSION=1.24.0

FROM golang:${GOLANG_VERSION}

ARG PROTOC_GEN_GO_VERSION=v1.31.0
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.3.0
ARG BUF_VERSION=v1.26.1
ARG PROTOC_VERSION=24.1

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION} && \
    go install github.com/bufbuild/buf/cmd/buf@${BUF_VERSION}

RUN apt update && apt install -y unzip && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /go

WORKDIR "/gen"
COPY proto .
COPY buf.gen.yaml .
RUN mkdir /.cache && chmod 0777 /.cache /gen

ENTRYPOINT ["/go/bin/buf", "generate"]
