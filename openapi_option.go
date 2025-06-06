package httpserver

import (
	"net/http"

	"github.com/ihezebin/openapi"
)

type OpenAPIOption func(*openapi.Route)

type OpenAPIOptions []OpenAPIOption

func NewOpenAPIOptions(opts ...OpenAPIOption) OpenAPIOptions {
	return opts
}

func WithOpenAPIDescription(description string) OpenAPIOption {
	return func(route *openapi.Route) {
		route.HasDescription(description)
	}
}

func WithOpenAPISummary(summary string) OpenAPIOption {
	return func(route *openapi.Route) {
		route.HasSummary(summary)
	}
}

func WithOpenAPIDeprecated() OpenAPIOption {
	return func(route *openapi.Route) {
		route.HasDeprecated(true)
	}
}

func WithOpenAPIResponseHeader(name string, param openapi.HeaderParam) OpenAPIOption {
	return func(route *openapi.Route) {
		route.HasResponseHeader(http.StatusOK, name, param)
	}
}

func mergeOpenAPIOptions(route *openapi.Route, options ...OpenAPIOption) *openapi.Route {
	for _, option := range options {
		option(route)
	}
	return route
}
