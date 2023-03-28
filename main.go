package main

import (
	"context"
	"fmt"
	"gg/app/controller"
	"gg/app/dao"
	"gg/app/interceptor"
	v1 "gg/app/proto/go/v1"
	"gg/config"
	"github.com/golang/glog"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// RegisterHttp 服务注册方法
type RegisterHttp func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

func main() {
	// 解析配置文件
	conf := config.GetConf()
	// 初始化数据库
	dao.Init(conf.Mysql)
	// grpc的中间件
	interceptors := []grpc.UnaryServerInterceptor{
		interceptor.Access(),
		interceptor.Recovery(),
		interceptor.SetHeader(),
	}
	// grpc的服务配置
	grpcOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(interceptors...)),
	}
	// 初始化一个grpcServer
	grpcServer := grpc.NewServer(grpcOptions...)
	// 注册服务
	v1.RegisterUserServiceServer(grpcServer, controller.MakeUser())
	// 启动grpc服务
	go grpcStart(grpcServer, conf)
	// 注册服务
	registerHttp := []RegisterHttp{
		v1.RegisterUserServiceHandlerFromEndpoint,
	}
	// 启动http服务
	go httpStart(conf, registerHttp)
	glog.Infof("服务启动成功, httpAddr: %s, grpcAddr: %s \n", conf.Addr.HttpAddr, conf.Addr.GrpcAddr)
	// 处理系统退出信号
	handleExitSignal(grpcServer)
}

// 启动http服务
func httpStart(conf *config.Config, registerHttp []RegisterHttp) {
	muxOptions := []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			lowKey := strings.ToLower(key)
			switch lowKey {
			case "x-request-id", "my-header":
				return lowKey, true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
		runtime.WithErrorHandler(interceptor.SimpleGrpcErrorHandler),
	}
	mux := runtime.NewServeMux(muxOptions...)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, register := range registerHttp {
		err := register(context.Background(), mux, conf.Addr.GrpcAddr, opts)
		if err != nil {
			panic(fmt.Sprintln("http服务启动失败, err:", err))
		}
	}
	// 注册中间件
	err := http.ListenAndServe(conf.Addr.HttpAddr, mux)
	if err != nil {
		panic(fmt.Sprintln("http服务启动失败, err:", err))
	}
}

// 启动grpc服务
func grpcStart(grpcServer *grpc.Server, conf *config.Config) {
	// 监听端口
	lis, err := net.Listen("tcp", conf.Addr.GrpcAddr)
	if err != nil {
		panic(fmt.Sprintln("grpc服务监听失败, addr:", conf.Addr.GrpcAddr, ", err:", err))
	}
	// grpc服务启动
	if err := grpcServer.Serve(lis); err != nil {
		panic(fmt.Sprintln("grpc服务启动失败, err:", err))
	}
}

// 处理系统退出信号
func handleExitSignal(grpc *grpc.Server) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	// 优雅退出
	grpc.GracefulStop()
	os.Exit(0)
}
