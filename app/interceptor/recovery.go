package interceptor

import (
	"context"
	"fmt"
	"gg/app/utils"
	"github.com/golang/glog"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery() grpc.UnaryServerInterceptor {
	return grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandlerContext(func(ctx context.Context,
		p interface{}) (err error) {
		glog.Errorf("panic triggered: %v \n", p)
		return status.Error(codes.Unknown, fmt.Sprintf("[%v] unknown server error", utils.GetRequestId(ctx)))
	}))
}
