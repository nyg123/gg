package interceptor

import (
	"context"
	"encoding/json"
	"errors"
	v1 "gg/app/proto/go/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

// SimpleGrpcErrorHandler 自定义grpc-gateway错误处理器 用于将grpc-gateway返回的错误转为http响应 并返回给客户端
// 在启动http服务时将此函数作为中间件传入
func SimpleGrpcErrorHandler(ctx context.Context,
	mux *runtime.ServeMux,
	marshal runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error) {
	var statusCode int

	// 尝试转换为HTTPStatusError并从中获取响应状态码和error
	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
		statusCode = customStatus.HTTPStatus
	}

	// 将 grpc headers 添加至 http response headers
	md, _ := runtime.ServerMetadataFromContext(ctx)
	for k, vs := range md.HeaderMD {
		for _, v := range vs {
			w.Header().Add(runtime.MetadataHeaderPrefix+k, v)
		}
	}

	// 从 grpc headers 中获取自定义响应状态码
	if statusCode == 0 {
		if values := md.HeaderMD.Get("x-http-code"); len(values) > 0 {
			if code, _ := strconv.Atoi(values[0]); code > 0 {
				statusCode = code
			}
			md.HeaderMD.Delete("x-http-code")
			w.Header().Del("Grpc-Metadata-X-Http-Code")
		}
	}

	// 将错误转为 grpc status 并从中获取 错误码 & 错误描述 & 响应状态码
	s := status.Convert(err)
	errorCode := int(s.Code())
	errorMsg := s.Message()
	if statusCode == 0 {
		statusCode = runtime.HTTPStatusFromCode(s.Code())
	}

	// 用框架规整的格式返回错误给客户端
	data := &v1.Response{
		Code:    int32(errorCode),
		Message: errorMsg,
	}
	bytes, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		log.Errorf(ctx, "返回http错误给客户端失败, err:", err)
	}
}
