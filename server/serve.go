package main

import (
	pb "../pb/google/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
 "fmt"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/metadata"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"grpcDemo/middleware"

	"sync"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService ...
var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloHttpRequest) (*pb.HelloHttpReply, error) {
	resp := new(pb.HelloHttpReply)
	md, _ := metadata.FromIncomingContext(ctx)
	sdkLang, _ := md["sdk-lang"]
	fmt.Println(sdkLang)
	resp.Message = "Hello " + in.Name + "."

	return resp, nil
}

func main() {

	Server()
}

func SslServer(){

	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(`/home/wendong/mygopath/src/grpcDemo/keys/server.pem`,
		`/home/wendong/mygopath/src/grpcDemo/keys/server.key`)
	if err != nil {
		fmt.Println("Failed to generate credentials %v", err)
		return
	}
	// 实例化grpc tls Server
	s := grpc.NewServer(grpc.Creds(creds))

	// 注册HelloService
	pb.RegisterHelloHttpServer(s, HelloService)

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}

func Server(){

	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}

	tM:=sync.Map{}
	// 实例化grpc tls Server
	s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(middleware.RateLimitMiddlewareFunc(middleware.GrpcRateLimit),
		middleware.UnaryServerPreInterceptor(&tM),
		middleware.UnaryServerAfterInterceptor(&tM),
		))

	// 注册HelloService
	pb.RegisterHelloHttpServer(s, HelloService)

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}