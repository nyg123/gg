package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var headerKey = "_header"

// Header 这里配置需要接收传递的header
type Header struct {
	MyHeader string `json:"my-header"`
}

func SetHeader() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		header := Header{}
		for k, v := range md {
			switch k {
			case "my-header":
				header.MyHeader = v[0]
			}
		}
		ctx = context.WithValue(ctx, headerKey, &header)
		return handler(ctx, req)
	}
}

func GetHeader(ctx context.Context) *Header {
	header, _ := ctx.Value(headerKey).(*Header)
	return header
}
