//go:generate protoc --go_out=plugins=grpc:. api.proto

/*
The api package defines api calls between the server and the client. It should
be included by go clients that want to call the server.
*/

package api

import (
	"errors"
	"google.golang.org/grpc"
)

type AutoClient struct {
	Addr   string
	conn   *grpc.ClientConn
	client ElTeeClient
}

func NewAutoClient(addr string) *AutoClient {
	return &AutoClient{
		Addr: addr,
	}
}

func (ac *AutoClient) Client() (ElTeeClient, error) {
	if ac == nil {
		return nil, errors.New("Uninitialized")
	}

	if ac.client != nil {
		return ac.client, nil
	}

	// It's empty, so try create it / connect it / whatever
	conn, err := grpc.Dial(ac.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ac.conn = conn
	ac.client = NewElTeeClient(conn)

	return ac.client, nil
}

func (ac *AutoClient) Close() {
	if ac == nil {
		return
	}

	if ac.conn != nil {
		ac.conn.Close()
		ac.conn = nil
	}

	ac.client = nil
}
