FROM golang:1.23 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/github.com/0x5d/wproc/
COPY src/ .
RUN go mod download
RUN go build -o wproc ./cmd/wproc

# Build
COPY . ./

ENTRYPOINT ["./wproc"]
