package client

import (
	"context"
	"time"

	"github.com/pugdag/pugdagd/cmd/pugdagwallet/daemon/server"

	"github.com/pkg/errors"

	"github.com/pugdag/pugdagd/cmd/pugdagwallet/daemon/pb"
	"google.golang.org/grpc"
)

// Connect connects to the pugdagwalletd server, and returns the client instance
func Connect(address string) (pb.KaspawalletdClient, func(), error) {
	// Connection is local, so 1 second timeout is sufficient
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(server.MaxDaemonSendMsgSize)))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, errors.New("pugdagwallet daemon is not running, start it with `pugdagwallet start-daemon`")
		}
		return nil, nil, err
	}

	return pb.NewKaspawalletdClient(conn), func() {
		conn.Close()
	}, nil
}
