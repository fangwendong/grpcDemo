package middleware

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"time"
	"sync"
	"errors"
	"fmt"
	"encoding/json"
)

func UnaryServerAfterInterceptor(dataTotalM *sync.Map) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}
		ip,err:=getClientIp(ctx)
		if err!=nil{
			return nil,errors.New("获取ip失败！")
		}

		var l int
		b,_:=json.Marshal(resp)
		l =len(b)
		fmt.Println("后置拦截器:", resp,"信息长度:",l, time.Now())
		dti, ok := dataTotalM.Load(ip)
		if ok {
			dt, ok := dti.(*dataDate)
			if ok {
				y, m, d := time.Now().Date()
				if (y*10000 + 100*int(m) + d) == dt.date {
					dt.total = dt.total + l
				} else {
					dt = &dataDate{date: y*10000 + 100*int(m) + d}
				}
				fmt.Println("total:",dt.total,"date:",dt.date)
				dataTotalM.Store(ip, dt)
			}
		} else {
			y, m, d := time.Now().Date()
			dt := &dataDate{date: y*10000 + 100*int(m) + d,total:l}
			dataTotalM.Store(ip, dt)
		}
		return resp, err
	}
}

