package protobuf

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingServerInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now().UnixMilli()
	log := info.Server.(*Server).logger
	requestorInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "cant define requestor ip")
	}
	resp, err := handler(ctx, req)
	end := time.Now().UnixMilli()
	log.Info(
		"%v [%v] %v %v %v %v %v",
		requestorInfo.Addr, time.Now().UTC().Format("02/Jan/2006:15:04:05 -0700"), info.FullMethod,
		"GRPC CALL", err, end-start, "protobuf-client",
	)
	return resp, err
}
