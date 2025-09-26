//go:build grpcgen
// +build grpcgen

package transfers

import (
	"context"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	transfersv1 "github.com/vaultone/api/internal/genproto/transfers/v1"
)

type Client struct {
	conn   *grpc.ClientConn
	client transfersv1.TransfersServiceClient
}

type Config struct {
	Addr        string
	DialTimeout time.Duration
	CallTimeout time.Duration
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	dialCtx, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(
		dialCtx,
		cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.DefaultConfig, MinConnectTimeout: 1 * time.Second}),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: transfersv1.NewTransfersServiceClient(conn)}, nil
}

func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) RequestTransfer(ctx context.Context, req *transfersv1.RequestTransferRequest, timeout time.Duration) (*transfersv1.RequestTransferResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.client.RequestTransfer(callCtx, req)
}
