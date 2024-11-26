package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authenticate is a grpc interceptor to do authentication on incoming request
func Authenticate(passphrase string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		// extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}

		// get authorization
		authorizationArr := md["authorization"]
		if len(authorizationArr) == 0 {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		// validate authorization
		authorization := strings.TrimPrefix(authorizationArr[0], "Bearer ")
		if authorization != passphrase {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		// continue request when ok
		return handler(ctx, req)
	}
}
