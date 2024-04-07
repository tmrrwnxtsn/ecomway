package integration

import pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"

type Client struct {
	client pb.IntegrationServiceClient
}

func NewClient(client pb.IntegrationServiceClient) *Client {
	return &Client{
		client: client,
	}
}
