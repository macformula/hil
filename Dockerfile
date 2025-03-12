FROM golang:1.24.1-bookworm

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /app/cmd/hilapp
RUN go build

WORKDIR /app/cmd/hilapp

ENTRYPOINT [ "/bin/bash" ]