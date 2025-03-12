FROM golang:1.24.1-bookworm

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
WORKDIR /app/cmd/hilapp
ENTRYPOINT [ "go", "build", "-buildvcs=false" ]