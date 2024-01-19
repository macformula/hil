package client

import (
	"context"

	"github.com/google/uuid"
	proto "github.com/macformula/hil/results/client/generated"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ResultsClient struct {
	addr                string
	conn                *grpc.ClientConn
	client              proto.TagTunnelClient
	pushReportsToGithub bool
}

func NewResultsClient(ip, port string, pushReportsToGithub bool) *ResultsClient {
	return &ResultsClient{
		addr:                ip + ":" + port,
		pushReportsToGithub: pushReportsToGithub,
	}
}

func (r *ResultsClient) Open(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, r.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "dial context")
	}

	r.conn = conn
	r.client = proto.NewTagTunnelClient(conn)

	return nil
}

func (r *ResultsClient) SubmitTag(ctx context.Context, tag string, value any) (bool, error) {
	request, err := createRequest(tag, value)
	if err != nil {
		return false, errors.Wrap(err, "create request")
	}

	reply, err := r.client.SubmitTag(ctx, request)
	if err != nil {
		return reply.IsPassing, errors.Wrap(err, "submit tag")
	}

	if !reply.Success {
		return false, errors.New(reply.Error)
	}

	return reply.IsPassing, nil
}

func (r *ResultsClient) CompleteTest(ctx context.Context, testId uuid.UUID, sequenceName string) (bool, error) {
	reply, err := r.client.CompleteTest(ctx, &proto.CompleteTestRequest{
		TestId:             testId.String(),
		SequenceName:       sequenceName,
		PushReportToGithub: r.pushReportsToGithub,
	})

	if err != nil {
		return false, errors.Wrap(err, "complete test")
	}

	return reply.TestPassed, nil
}

func (r *ResultsClient) SubmitError(ctx context.Context, err error) error {
	_, submitErr := r.client.SubmitError(ctx, &proto.SubmitErrorRequest{Error: err.Error()})
	if submitErr != nil {
		return errors.Wrap(err, "submit error")
	}

	return nil
}

func createRequest(tag string, data any) (*proto.SubmitTagRequest, error) {
	request := &proto.SubmitTagRequest{Tag: tag}

	switch data.(type) {
	case int32:
		request.Data = &proto.SubmitTagRequest_ValueInt{ValueInt: data.(int32)}
	case float32:
		request.Data = &proto.SubmitTagRequest_ValueFloat{ValueFloat: data.(float32)}
	case string:
		request.Data = &proto.SubmitTagRequest_ValueStr{ValueStr: data.(string)}
	case bool:
		request.Data = &proto.SubmitTagRequest_ValueBool{ValueBool: data.(bool)}
	default:
		return nil, errors.Errorf("unsupported data type for tag submission (%T)", data)
	}

	return request, nil
}
