package grpcx

import (
	"context"
	"net/http"
	"strings"

	"github.com/fagongzi/log"
	"github.com/labstack/echo"
	md "github.com/labstack/echo/middleware"
)

type httpServer struct {
	addr   string
	server *echo.Echo
}

func newHTTPServer(addr string) *httpServer {
	return &httpServer{
		addr:   addr,
		server: echo.New(),
	}
}

func (s *httpServer) start() error {
	s.server.Use(md.Recover())

	log.Infof("rpc: start a grpc http proxy server at %s", s.addr)
	return s.server.Start(s.addr)
}

func (s *httpServer) stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *httpServer) addService(service Service) {
	if service.Metadata.HTTP != nil {
		m := strings.ToUpper(service.Metadata.HTTP.Method)
		switch m {
		case echo.GET:
			s.server.GET(service.Metadata.HTTP.Path, func(c echo.Context) error {
				return s.handleHTTP(c, service)
			})
		case echo.PUT:
			s.server.PUT(service.Metadata.HTTP.Path, func(c echo.Context) error {
				return s.handleHTTP(c, service)
			})
		case echo.DELETE:
			s.server.DELETE(service.Metadata.HTTP.Path, func(c echo.Context) error {
				return s.handleHTTP(c, service)
			})
		case echo.POST:
			s.server.POST(service.Metadata.HTTP.Path, func(c echo.Context) error {
				return s.handleHTTP(c, service)
			})
		}
	}
}

func (s *httpServer) handleHTTP(c echo.Context, service Service) error {
	if service.invoker == nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	data, err := service.invoker(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, data)
}
