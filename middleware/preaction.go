package middleware

import (
	"errors"
	"sync"

	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"fmt"
)

const numMax = 100

type dataDate struct {
	total int
	date  int
}

func UnaryServerPreInterceptor(dataTotalM *sync.Map) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ip,err:=getClientIp(ctx)
		if err!=nil{
			return nil,errors.New("获取ip失败！")
		}
		dti, ok := dataTotalM.Load(ip)
		if ok {
			dt, ok := dti.(*dataDate)
			if ok {
				fmt.Println("[pre] total:",dt.total,"date:",dt.date)
				y, m, d := time.Now().Date()
				if (y*10000 + 100*int(m) + d) == dt.date {
					if dt.total >= numMax {
						return nil, errors.New("today data sum over !")
					}
				}
			}
		}
		return handler(ctx, req)
	}
}
