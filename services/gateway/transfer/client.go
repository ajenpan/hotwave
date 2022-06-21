package transfer

import (
	protobuf "google.golang.org/protobuf/proto"
)

type Client struct {
}

func NewClient() *Client {
	ret := &Client{}

	// conn, err := grpc.Dial("unix:///var/lib/test.socket", grpc.WithInsecure())
	// if err!=nil{
	// 	return nil
	// }
	// client := protocol.NewGatewayClient(conn)
	// SendMessageToUse
	// proto.SendMessageToUserRequest
	return ret
}

func (c *Client) SendMessageToUse(nodeId string, userId int64, message protobuf.Message) error {

	// SendMessageToUserRequest
	// proto.SendMessageToUserRequest
	return nil
}
