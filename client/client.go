package main

import (
	"google.golang.org/grpc"
	pb "../pb/google/api"
	"context"
	"fmt"

	"grpcDemo/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	"time"

)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	// 连接
	//cred,err:=credentials.NewClientTLSFromFile(`/home/wendong/mygopath/src/grpcDemo/keys/server.pem`,"sf")
	//if err != nil {
	//	fmt.Println("NewClientTLSFromFile fail!",err)
	//	return
	//}
	//var ch chan int
	conn, err := grpc.Dial(Address,
		grpc.WithInsecure(),
	grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.OperationMiddlewareFunc(middleware.SetContext,map[string][]string{"sdk-lang": {"r"}},
		),
		),
	))
	if err != nil {
		fmt.Println("dial fail!", err)
		return
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloHttpClient(conn)

	// 调用方法
	reqBody := new(pb.HelloHttpRequest)
	reqBody.Name = "gRPC"
	now:=time.Now()

	//看看1s内是否只通过2个请求

		go func() {
			for{
				r, err := c.SayHello(context.Background(), reqBody)
				if err == nil {
					fmt.Println(r.Message,"耗时:",time.Since(now))
					now=time.Now()
				}
			}

		}()
	time.Sleep(time.Second*8)
}

/*
测试结果： 4s内通过了8次
Hello gRPC. 耗时: 1.982095ms
Hello gRPC. 耗时: 184.766µs
Hello gRPC. 耗时: 999.66874ms
Hello gRPC. 耗时: 374.533µs
Hello gRPC. 耗时: 999.663983ms
Hello gRPC. 耗时: 368.11µs
Hello gRPC. 耗时: 999.653043ms
Hello gRPC. 耗时: 388.048µs

 */