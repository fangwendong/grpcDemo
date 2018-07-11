package main

import (
	"flag"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
	"github.com/golang/glog"
	gw "../pb/google/api"
	"google.golang.org/grpc/credentials"
	"fmt"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:50052", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// 连接
	cred,err:=credentials.NewClientTLSFromFile(`/home/wendong/mygopath/src/grpcDemo/keys/server.pem`,
		"sf")
	if err != nil {
		fmt.Println("NewClientTLSFromFile fail!",err)
		return err
	}
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(cred)}
	err = gw.RegisterHelloHttpHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts )
	if err != nil {
		return err
	}

	return  http.ListenAndServeTLS(":8888", `/home/wendong/mygopath/src/grpcDemo/keys/server.pem`,
		`/home/wendong/mygopath/src/grpcDemo/keys/server.key`, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
