package frame

import (
	"context"

	"google.golang.org/grpc"
	pb "google.golang.org/protobuf/proto"

	"hotwave/frame/proto"
)

type Client struct {
}

func (c *Client) Send(ctx context.Context, msg pb.Message) error {

	conn, err := grpc.Dial("")
	if err != nil {
		return err
	}
	client := proto.NewNodeBaseClient(conn)
	steam, err := client.UserMessage(ctx)
	if err != nil {
		return err
	}
	steam.Send(&proto.UserMessageWraper{})
	return nil
}
