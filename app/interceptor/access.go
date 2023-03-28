package interceptor

import (
	"context"
	"encoding/json"
	"gg/app/utils"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"strconv"
	"strings"
	"time"
)

func Access() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		// 尝试从header中获取requestId
		md, ok := metadata.FromIncomingContext(ctx)
		var traceId string
		if ok {
			if v, ok := md[utils.RequestIDKey]; ok {
				// 如果传入了x-request-id，就使用传入的
				traceId = v[0]
			} else {
				// 如果没有传入x-request-id就自己生成一个
				traceId = strconv.FormatInt(utils.GetSnowflake(), 10)
			}
		}

		// 设置 后续发起的 grpc 请求的 metadata: x-forwarded-for 值
		p, _ := peer.FromContext(ctx)
		h, _, _ := net.SplitHostPort(p.Addr.String())
		f, ok := md["x-forwarded-for"]
		f = append(f, h)

		// 设置 context trace id & 设置 grpc 请求/响应头
		ctx = context.WithValue(ctx, utils.RequestIDKey, traceId)
		ctx = metadata.NewOutgoingContext(ctx,
			metadata.Pairs("X-Request-Id", traceId, "X-Forwarded-For", strings.Join(f, ", ")))
		_ = grpc.SetHeader(ctx, metadata.Pairs("X-Request-Id", traceId))

		resp, err = handler(ctx, req)

		// 获取响应状态码
		var code int
		if st, ok := grpc.ServerTransportStreamFromContext(ctx).(interface {
			Header() (metadata.MD, error)
		}); ok {
			md, _ := st.Header()
			if values := md.Get("x-http-code"); len(values) > 0 {
				if c, _ := strconv.Atoi(values[0]); c > 0 {
					code = c
				}
			}
		}
		headers := ""
		for key, val := range md {
			headers += key + ": " + strings.Join(val, " ") + " \n"
		}
		body, _ := json.Marshal(req)
		glog.Infof("[access] [%s] [%d] [%dms] GRPC %s \n%v \n%s \n---------- \n", h, code,
			time.Since(startTime).Milliseconds(), info.FullMethod, headers, body)
		return resp, err
	}
}
