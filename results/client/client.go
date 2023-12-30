package client

import (
	"context"
	"fmt"
	proto "github.com/macformula/hil/results/client/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const addr = "localhost:8080"

// Data generic to simplify submitting tags and only allow for expected types.
// ~ syntax allows for us to pass other types whose underlying type is one of these primitive types.
// i.e. type Xyz int32, not sure if we'd need it, but it could be useful.
type Data interface {
	~float32 | ~int32 | ~bool | ~string
}

func SubmitTag[T Data](tag string, data T) (bool, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}

	c := proto.NewTagTunnelClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := c.SubmitTag(ctx, createRequest(tag, data))
	if err != nil {
		return reply.Success, err
	}

	closeConn(conn)
	return reply.Success, nil
}

func CompleteTest() (bool, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}

	c := proto.NewTagTunnelClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := c.CompleteTest(ctx, &emptypb.Empty{})
	if err != nil {
		return reply.Success, err
	}

	closeConn(conn)
	return reply.Success, nil
}

func createRequest[T Data](tag string, data T) *proto.SubmitTagRequest {
	typeStr := fmt.Sprintf("%T", data)
	request := &proto.SubmitTagRequest{Tag: tag}
	switch typeStr {
	case "int32":
		request.Data = &proto.SubmitTagRequest_ValueInt{ValueInt: any(data).(int32)}
	case "float32":
		request.Data = &proto.SubmitTagRequest_ValueFloat{ValueFloat: any(data).(float32)}
	case "string":
		request.Data = &proto.SubmitTagRequest_ValueStr{ValueStr: any(data).(string)}
	case "bool":
		request.Data = &proto.SubmitTagRequest_ValueBool{ValueBool: any(data).(bool)}
	}
	return request
}

func closeConn(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		zap.L().Error("encountered an error when closing gRPC connection", zap.Error(err))
	}
}
