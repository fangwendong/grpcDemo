package middleware

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type OperationMiddleware func(ctx context.Context,data map[string][]string) (context.Context, error)

func OperationMiddlewareFunc(middlewareFunc OperationMiddleware,data map[string][]string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,opts ...grpc.CallOption) ( err error) {
		newCtx, err := middlewareFunc(ctx,data)
		if err != nil {
			return err
		}
		return invoker(newCtx, method,req, reply ,cc,opts...)
	}
}

func SetContext(ctx context.Context,data map[string][]string) (newCtx context.Context, err error) {
	md, _:= metadata.FromIncomingContext(ctx)
	//md :=map[string][]string{}
	if len(md)==0{
		md = map[string][]string{}
	}

	for k,v:=range data{
		md.Set(k,v...)
	}
	newCtx = metadata.NewOutgoingContext(ctx,md)
		return newCtx, nil
}
