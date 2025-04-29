package httpserver

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ihezebin/openapi"
	"go.opentelemetry.io/otel/trace"

	"github.com/ihezebin/olympus/httpserver/middleware"
	"github.com/ihezebin/olympus/logger"
)

var ctx = context.Background()

func TestServer(t *testing.T) {
	server := NewServer(
		WithPort(8000),
		WithServiceName("test_server"),
		WithMiddlewares(middleware.LoggingRequestWithoutHeader(), middleware.LoggingResponse()),
	)

	server.RegisterRoutes(&HelloRouter{})

	err := server.RegisterOpenAPIUI("/stoplight", StoplightUI)
	if err != nil {
		t.Fatal(err)
	}
	err = server.RegisterOpenAPIUI("/swagger", SwaggerUI)
	if err != nil {
		t.Fatal(err)
	}
	err = server.RegisterOpenAPIUI("/redoc", RedocUI)
	if err != nil {
		t.Fatal(err)
	}
	err = server.RegisterOpenAPIUI("/rapidoc", RapidocUI)
	if err != nil {
		t.Fatal(err)
	}

	if err := server.RunWithNotifySignal(ctx); err != nil {
		t.Fatal(err)
	}
}

/*
curl --location 'http://127.0.0.1:8000/hello/ping' \
--header 'Traceparent: 00-5e64f77760384d153783c96049550881-b7ad6b7169203331-01'

这个 Traceparent 的格式：
00-5e64f77760384d153783c96049550881-b7ad6b7169203331-01
00-traceId-spanId-flags

traceId: 5e64f77760384d153783c96049550881
spanId: b7ad6b7169203331
flags: 01

相关文档：https://opentelemetry.io/docs/specs/otel/context/api-propagators/#w3c-trace-context-requirements
*/
func TestServerWithOtel(t *testing.T) {
	server := NewServer(
		WithPort(8000),
		WithServiceName("test_server"),
		WithOtel(true),
	)

	server.RegisterRoutes(&HelloRouter{})

	if err := server.RunWithNotifySignal(ctx); err != nil {
		t.Fatal(err)
	}
}

type HelloRouter struct {
}

func (h *HelloRouter) RegisterRoutes(router Router) {
	group := router.Group("/hello")
	group.POST("/world", NewHandler(h.Hello))
	group.GET("/ping", NewHandler(h.Ping), WithOpenAPIDeprecated(), WithOpenAPIResponseHeader("Token", openapi.HeaderParam{
		Description: "认证 JWT",
	}))
}

type HelloReq struct {
	Content string `json:"content" form:"content"`
	Id      string `json:"id" form:"id"`
}

// HelloResp 测试的 hello 响应数据
type HelloResp struct {
	Message string `json:"message"`
}

func (h *HelloRouter) Hello(c *gin.Context, req *HelloReq) (resp *HelloResp, err error) {
	return &HelloResp{Message: req.Content}, nil
}

func (h *HelloRouter) Ping(c *gin.Context, req map[string]interface{}) (resp string, err error) {
	ctx := c.Request.Context()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	logger.Infof(ctx, "traceId: %s, req: %+v", traceId, req)
	return "pong!" + traceId, nil
}
