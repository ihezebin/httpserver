package httpserver

import "github.com/gin-gonic/gin"

type RouterOptions struct {
	PreMiddlewares  []gin.HandlerFunc
	PostMiddlewares []gin.HandlerFunc
	OpenAPIOptions  OpenAPIOptions
	PathRegister    func(path string)
}

type RouterOption func(*RouterOptions)

func mergeRouterOptions(options ...RouterOption) *RouterOptions {
	routerOptions := &RouterOptions{}
	for _, option := range options {
		option(routerOptions)
	}
	return routerOptions
}

func WithPreMiddlewares(middlewares ...gin.HandlerFunc) RouterOption {
	return func(options *RouterOptions) {
		options.PreMiddlewares = append(options.PreMiddlewares, middlewares...)
	}
}

func WithPostMiddlewares(middlewares ...gin.HandlerFunc) RouterOption {
	return func(options *RouterOptions) {
		options.PostMiddlewares = append(options.PostMiddlewares, middlewares...)
	}
}

func WithOpenAPIOptions(opts ...OpenAPIOption) RouterOption {
	return func(options *RouterOptions) {
		options.OpenAPIOptions = NewOpenAPIOptions(opts...)
	}
}

func WithPathRegister(pathRegister func(path string)) RouterOption {
	return func(options *RouterOptions) {
		options.PathRegister = pathRegister
	}
}
