package main

import (
	proto "client/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"time"
)

const addr = "localhost:8080"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewTagTunnelClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	value := rand.Float32()
	_, err = c.SubmitTag(ctx, &proto.SubmitTagRequest{
		Tag:  "xyz",
		Data: &proto.SubmitTagRequest_ValueFloat{ValueFloat: value},
	})
	fmt.Printf("Tag submitted with value: %f", value)
}
