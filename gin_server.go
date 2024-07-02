package ji

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	addr string // 启动ip:port，空则默认是 localhost:8080
	c    chan error

	engine *gin.Engine // 如 gin.Default()
}

func NewGinServer(addr string, engine *gin.Engine) *GinServer {
	server := &GinServer{
		addr:   addr,
		c:      make(chan error),
		engine: engine,
	}
	var _ Server = server
	return server
}

func (m *GinServer) Start(ctx context.Context) (err error) {
	m.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return m.engine.Run(m.addr)
}

func (m *GinServer) AsyncStart(ctx context.Context) {
	go func() {
		err := m.Start(ctx)
		if err != nil {
			fmt.Println("[AsyncStart] Start panic", err)
			m.c <- err
		}
	}()
}

func (m *GinServer) ErrorC() (c chan error) {
	return m.c
}
