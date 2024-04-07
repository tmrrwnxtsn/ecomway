package engine

import pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"

type Client struct {
	client pb.EngineServiceClient
}

func NewClient(client pb.EngineServiceClient) *Client {
	return &Client{
		client: client,
	}
}
