package ji

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-gonic/gin"
)

type Server interface {
	// Start 开始提供服务
	Start(ctx context.Context) (err error)
	// AsyncStart 异步启动
	AsyncStart(ctx context.Context)
	// ErrorC 异步启动情况下，用 channel 返回启动结果（error）
	ErrorC() (c chan error)
}

func StartServerByName(ctx context.Context, serviceName string, processor thrift.TProcessor, engine *gin.Engine) (err error) {
	serviceInfo, err := GetServiceInfoByName(serviceName)
	if err != nil {
		return
	}

	var servers []Server
	if processor != nil {
		thriftServer := NewThriftServer(serviceInfo.Protocol, serviceInfo.IsBuffered, serviceInfo.IsFramed,
			serviceInfo.GetAddr(), serviceInfo.UseSecure, processor)
		servers = append(servers, thriftServer)
	}

	if engine != nil {
		ginServer := NewGinServer(serviceInfo.GetWebAddr(), engine)
		servers = append(servers, ginServer)
	}

	c := make(chan error)
	for _, server := range servers {
		go func(server Server) {
			c <- server.Start(ctx)
		}(server)
	}

	if len(servers) == 0 {
		close(c)
		fmt.Println("no server start!")
		return
	}

	fmt.Println("server async started.")

	select {
	case err := <-c:
		fmt.Println("server start error", err)
		return err
	}
}
