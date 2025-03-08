FROM golang:1.24.1-bookworm

RUN apt install ca-certificates git

WORKDIR /app

RUN apt-get install can-utils
RUN modprobe vcan
RUN ip link add dev vcan0 type vcan && \
    ip link set vcan0 mtu 16 && \
    ip link set up vcan0
RUN ip link add dev vcan1 type vcan && \
    ip link set vcan1 mtu 16 && \
    ip link set up vcan1

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /app/cmd/hilapp
RUN go build

WORKDIR /

CMD ["/bin/sh", "-c", "./main; exec sh"]