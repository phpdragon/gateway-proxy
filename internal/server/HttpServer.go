package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	Server           *http.Server
	shutdownFinished chan struct{}
}

// ListenAndServe 监听端口并提供服务
func (s *HttpServer) ListenAndServe() (err error) {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	config.Logger().Infof("Listening on port %s", s.Server.Addr)
	config.Logger().Infof("Open http://localhost:%s in the browser", s.Server.Addr)

	err = s.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	} else if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
		return
	}

	config.Logger().Info("Waiting for http server shutdown finishing...")
	<-s.shutdownFinished
	config.Logger().Info("Http server shutdown finished")

	return
}

// WaitForExitingSignal 监听退出信号，实现安全退出
func (s *HttpServer) WaitForExitingSignal(timeout time.Duration) {
	waiter := make(chan os.Signal, 1) // 按文档指示，至少设置1的缓冲
	signal.Notify(waiter, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-waiter // 阻塞直到有指定信号传入

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// We received an interrupt signal, shut down.
	if err := s.Server.Shutdown(ctx); err != nil {
		config.Logger().Error("Http server shutdown error", err)
	}

	config.Logger().Info("Http server shutdown successfully")

	close(s.shutdownFinished)
}
