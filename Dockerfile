FROM golang:1.13.0 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /workspace

RUN GO111MODULE=on go build -o out .

FROM ubuntu:18.04
WORKDIR /
COPY --from=builder /workspace/out /usr/local/bin/github-api-auth-proxy
ENTRYPOINT ["/usr/local/bin/github-api-auth-proxy"]
