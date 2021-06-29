package pkg

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GRPClient struct {
	conn      *grpc.ClientConn
	dialOpts  []grpc.DialOption
	timeout   time.Duration
	tlsConfig *tls.Config
}

func NewGRPCClient(cfg *Config) (*GRPClient, error) {
	client := &GRPClient{}

	tlscfg, err := NewTLSConfig(cfg.TLS)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create tls config for grpc client")
	}
	client.tlsConfig = tlscfg

	client.dialOpts = append(client.dialOpts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    1 * time.Minute,
		Timeout: 20 * time.Second,
	}))
	client.dialOpts = append(client.dialOpts, grpc.WithBlock())
	client.dialOpts = append(client.dialOpts, grpc.FailOnNonTempDialError(true))

	client.timeout = cfg.Timeout

	return client, nil
}

func (client *GRPClient) Connect(ctx context.Context, address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, client.dialOpts...)

	if client.tlsConfig != nil {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(client.tlsConfig)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	c, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	log.Info("dialing...")
	conn, err := grpc.DialContext(c, address, opts...)
	if err != nil {
		log.WithError(err).Error("failed to create client connection")
		return nil, errors.WithMessagef(
			errors.WithStack(err),
			"failed to create new client connection to %s",
			address,
		)
	}

	log.Info("successfully established client connection")
	client.conn = conn
	return conn, nil
}
