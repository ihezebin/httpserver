package httpserver

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ServerOptions struct {
	Port            uint               `json:"port" yaml:"port" toml:"port"`
	Daemon          bool               `json:"daemon" yaml:"daemon" toml:"daemon"`
	Middlewares     []gin.HandlerFunc  `json:"middlewares" yaml:"middlewares" toml:"middlewares"`
	ServiceName     string             `json:"service_name" yaml:"service_name" toml:"service_name"`
	HiddenRoutesLog bool               `json:"hidden_routes_log" yaml:"hidden_routes_log" toml:"hidden_routes_log"`
	Pprof           bool               `json:"pprof" yaml:"pprof" toml:"pprof"`
	OpenAPInfo      *openapi3.Info     `json:"openap_info" yaml:"openap_info" toml:"openap_info"`
	OpenAPIServers  []openapi3.Server  `json:"openap_server" yaml:"openap_server" toml:"openap_server"`
	Metrics         bool               `json:"metrics" yaml:"metrics" toml:"metrics"`
	TraceExporter   trace.SpanExporter `json:"trace_exporter" yaml:"trace_exporter" toml:"trace_exporter"`
	LogProcessor    log.Processor      `json:"log_processor" yaml:"log_processor" toml:"log_processor"`
}

type ServerOption func(*ServerOptions)

func mergeServerOptions(opts ...ServerOption) *ServerOptions {
	opt := &ServerOptions{
		Port:    8080,
		Pprof:   true,
		Metrics: true,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// WithTraceExporter 设置 trace exporter
// otlptracehttp.New(ctx,
//
//	otlptracehttp.WithInsecure(),
//	otlptracehttp.WithEndpoint("localhost:54318"),
//
// )
func WithTraceExporter(exporter trace.SpanExporter) ServerOption {
	return func(o *ServerOptions) {
		o.TraceExporter = exporter
	}
}

// WithMetrics 使用 prometheus 作为 metrics 的 provider reader
func WithMetrics() ServerOption {
	return func(o *ServerOptions) {
		o.Metrics = true
	}
}

// WithLogProcessor 设置 log processor
// 日志需要实现通过 otellog.Logger Emit 日志内容
func WithLogProcessor(exporter log.Exporter, opts ...log.BatchProcessorOption) ServerOption {
	return func(o *ServerOptions) {
		o.LogProcessor = log.NewBatchProcessor(exporter, opts...)
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

func WithOpenAPIServer(server ...openapi3.Server) ServerOption {
	return func(o *ServerOptions) {
		o.OpenAPIServers = append(o.OpenAPIServers, server...)
	}
}
