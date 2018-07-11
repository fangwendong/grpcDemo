package middleware

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync"
	"time"
	"myquant.cn/platform/core/util/ratelimit"
)

var rateLimitM sync.Map  //用来存放ip对应的limiter
type RateLimitMiddleware func(ctx context.Context)(context.Context, error)
func RateLimitMiddlewareFunc(middlewareFunc RateLimitMiddleware)grpc.UnaryServerInterceptor{
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx, err := middlewareFunc(ctx)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}


//通过ip进行频率限制
func GrpcRateLimit(ctx context.Context)(context.Context,error){
	//1.获取客户端ip
	ip,err:=getClientIp(ctx)
	if err!=nil{
		return ctx,err
	}
	//fmt.Println("clientIp:",ip)
	//2.限制ip
	value,ok:=rateLimitM.Load(ip)
	if !ok{
		limiter:=ratelimit.NewWindowLimiter(time.Second*1,2)
		rateLimitM.Store(ip,limiter)
		return ctx,nil
	}
	l,_:=value.(ratelimit.Limiter)
	//now:=time.Now()
	_, ok = l.TakeWithTimestamp(1, time.Now())
	if !ok{
		//fmt.Println("限制流量了")
		//fmt.Println("频率限制阻塞时长:",time.Since(now),wait,ok)
		return ctx,errors.New("限制流量了,无法通过！")
	}

	return ctx,nil

	return ctx,nil
}

func getClientIp(ctx context.Context)(string,error){
	var ip string
	pr,ok:=peer.FromContext(ctx)
	if !ok{
		return ip,errors.New("FromContext(ctx) failed !")
	}

	if pr.Addr ==net.Addr(nil){
		return ip,errors.New("peer.Addr is nil !")
	}

	addSlice:=strings.Split(pr.Addr.String(),":")
	if len(addSlice)==0{
		return ip,errors.New("len(addSlice)==0")
	}

	return addSlice[0],nil
}