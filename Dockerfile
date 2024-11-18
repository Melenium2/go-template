FROM golang:1.23 as build

ARG WORKDIR=${GOPATH}/src/qleanlabs/app/
ARG GITLAB_GO_TOKEN

WORKDIR ${WORKDIR}

# build service
COPY . .
RUN go build -tags timetzdata -o /usr/bin/service-entrypoint ./cmd/service/

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=build /usr/bin/service-entrypoint /usr/bin/
COPY db /db

EXPOSE 4000
CMD [ "service-entrypoint" ]

