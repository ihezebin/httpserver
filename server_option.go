package httpserver

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ServerOptions struct {
	Port            uint                  `json:"port" yaml:"port" toml:"port"`
	Daemon          bool                  `json:"daemon" yaml:"daemon" toml:"daemon"`
	Middlewares     []gin.HandlerFunc     `json:"middlewares" yaml:"middlewares" toml:"middlewares"`
	ServiceName     string                `json:"service_name" yaml:"service_name" toml:"service_name"`
	HiddenRoutesLog bool                  `json:"hidden_routes_log" yaml:"hidden_routes_log" toml:"hidden_routes_log"`
	Metrics         bool                  `json:"metrics" yaml:"metrics" toml:"metrics"`
	Pprof           bool                  `json:"pprof" yaml:"pprof" toml:"pprof"`
	OpenAPInfo      *openapi3.Info        `json:"openap_info" yaml:"openap_info" toml:"openap_info"`
	OpenAPIServer   *openapi3.Server      `json:"openap_server" yaml:"openap_server" toml:"openap_server"`
	Otel            bool                  `json:"otel_trace" yaml:"otel_trace" toml:"otel_trace"`
	OtelInit        func() []ShutdownFunc `json:"otel_init" yaml:"otel_init" toml:"otel_init"`
}

type ServerOption func(*ServerOptions)

func mergeServerOptions(opts ...ServerOption) *ServerOptions {
	opt := &ServerOptions{
		Port: 8080,
		Otel: true,
		OtelInit: func() []ShutdownFunc {
			tp := trace.NewTracerProvider()
			otel.SetTracerProvider(tp)
			propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
			otel.SetTextMapPropagator(propagator)
			shutdowns := []ShutdownFunc{tp.Shutdown}
			return shutdowns
		},
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithOtelInit(init func() []ShutdownFunc) ServerOption {
	return func(o *ServerOptions) {
		o.OtelInit = init
	}
}

func WithOtel(enabled bool) ServerOption {
	return func(o *ServerOptions) {
		o.Otel = enabled
	}
}

func WithPort(port uint) ServerOption {
	return func(o *ServerOptions) {
		o.Port = port
	}
}

func WithDaemon(daemon bool) ServerOption {
	return func(o *ServerOptions) {
		o.Daemon = daemon
	}
}

func WithMiddlewares(middlewares ...gin.HandlerFunc) ServerOption {
	return func(o *ServerOptions) {
		o.Middlewares = middlewares
	}
}

func WithMetrics() ServerOption {
	return func(o *ServerOptions) {
		o.Metrics = true
	}
}

func WithPprof() ServerOption {
	return func(o *ServerOptions) {
		o.Pprof = true
	}
}

const (
	ServiceNameKey = "service"
)

func WithServiceName(name string) ServerOption {
	return func(o *ServerOptions) {
		o.ServiceName = name
	}
}

func WithHiddenRoutesLog() ServerOption {
	return func(o *ServerOptions) {
		o.HiddenRoutesLog = true
	}
}

func WithOpenAPInfo(info openapi3.Info) ServerOption {
	return func(o *ServerOptions) {
		o.OpenAPInfo = &info
	}
}

func WithOpenAPIServer(server openapi3.Server) ServerOption {
	return func(o *ServerOptions) {
		o.OpenAPIServer = &server
	}
}
