package httpserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ihezebin/openapi"
)

type Router interface {
	Group(...string) Router
	Use(...gin.HandlerFunc) Router
	Any(string, handlerGenerator, ...OpenAPIOption)
	AnyWithOptions(string, handlerGenerator, ...RouterOption)
	GET(string, handlerGenerator, ...OpenAPIOption)
	GetWithOptions(string, handlerGenerator, ...RouterOption)
	POST(string, handlerGenerator, ...OpenAPIOption)
	PostWithOptions(string, handlerGenerator, ...RouterOption)
	DELETE(string, handlerGenerator, ...OpenAPIOption)
	DeleteWithOptions(string, handlerGenerator, ...RouterOption)
	PATCH(string, handlerGenerator, ...OpenAPIOption)
	PatchWithOptions(string, handlerGenerator, ...RouterOption)
	PUT(string, handlerGenerator, ...OpenAPIOption)
	PutWithOptions(string, handlerGenerator, ...RouterOption)
	OPTIONS(string, handlerGenerator, ...OpenAPIOption)
	OptionsWithOptions(string, handlerGenerator, ...RouterOption)
	HEAD(string, handlerGenerator, ...OpenAPIOption)
	HeadWithOptions(string, handlerGenerator, ...RouterOption)
}

type openapiRouter struct {
	prefix    string
	ginRouter gin.IRouter
	openapi   *openapi.API
}

func (r *openapiRouter) Group(prefixes ...string) Router {
	prefix := strings.Join(prefixes, "/")

	return &openapiRouter{
		prefix:    strings.ReplaceAll(strings.Join([]string{r.prefix, prefix}, "/"), "//", "/"),
		ginRouter: r.ginRouter.Group(prefix),
		openapi:   r.openapi,
	}
}

func (r *openapiRouter) mergePath(paths ...string) string {
	path := strings.ReplaceAll(strings.Join(paths, "/"), "//", "/")
	if strings.Contains(path, "//") {
		return r.mergePath(paths...)
	}

	return path
}

func (r *openapiRouter) Use(middleware ...gin.HandlerFunc) Router {
	r.ginRouter.Use(middleware...)
	return r
}

func (r *openapiRouter) handle(method string, path string, h handlerGenerator, routerOptions *RouterOptions) {
	requestBody, responseBody, query, params, requestHeader, responseHeader, handlerFunc := h()

	ginFuncs := make([]gin.HandlerFunc, 0, 1)
	ginFuncs = append(ginFuncs, routerOptions.PreMiddlewares...)
	ginFuncs = append(ginFuncs, handlerFunc)
	ginFuncs = append(ginFuncs, routerOptions.PostMiddlewares...)

	// register gin route
	r.ginRouter.Handle(method, path, ginFuncs...)

	// handle openapi path
	path = r.mergePath(r.prefix, path)
	// handle openapi route
	route := r.openapi.Route(method, path)
	route = mergeOpenAPIOptions(route, routerOptions.OpenAPIOptions...)
	operationID := strings.ReplaceAll(path, "/", "_")
	operationID = strings.TrimLeft(operationID, "_")
	operationID = strings.TrimRight(operationID, "_")

	route.HasOperationID(operationID)

	if requestBody != nil {
		route.HasRequestModel(*requestBody)
	}

	route.HasResponseModel(http.StatusInternalServerError, openapi.ModelOf[Body[EmptyType]]())
	if responseBody != nil {
		route.HasResponseModel(http.StatusOK, *responseBody)
	}

	if len(query) > 0 {
		for k, v := range query {
			route.HasQueryParameter(k, v)
		}
	}

	// path 里面有，但是由于 uri tag 添加的 param 要删除
	realExistParam := make(map[string]bool)
	// path 里面包含 :id 格式的，添加 param
	if strings.Contains(path, ":") {
		for _, param := range strings.Split(path, "/") {
			if strings.Contains(param, ":") {
				realExistParam[param] = true
				param = strings.TrimLeft(param, ":")
				pathParam, ok := params[param]
				if !ok {
					pathParam = openapi.PathParam{
						Description: "",
						Type:        openapi.PrimitiveTypeString,
					}
				}
				params[param] = pathParam
			}
		}
	}

	if len(params) > 0 {
		for k, v := range params {
			if realExistParam[k] {
				route.HasPathParameter(k, v)
			}
		}
	}

	if len(requestHeader) > 0 {
		for k, v := range requestHeader {
			route.HasHeaderParameter(k, v)
		}
	}

	if len(responseHeader) > 0 {
		for k, v := range responseHeader {
			route.HasResponseHeader(http.StatusOK, k, v)
		}
	}
}

func (r *openapiRouter) Any(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.AnyWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) AnyWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	anyMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}

	routerOptions := mergeRouterOptions(options...)

	for _, method := range anyMethods {
		r.handle(method, path, h, routerOptions)
	}
}

func (r *openapiRouter) GET(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.GetWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) GetWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)

	r.handle(http.MethodGet, path, h, routerOptions)
}

func (r *openapiRouter) POST(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.PostWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) PostWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)
	r.handle(http.MethodPost, path, h, routerOptions)
}

func (r *openapiRouter) DELETE(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.DeleteWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) DeleteWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)
	r.handle(http.MethodDelete, path, h, routerOptions)
}

func (r *openapiRouter) PATCH(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.PatchWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) PatchWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)
	r.handle(http.MethodPatch, path, h, routerOptions)
}

func (r *openapiRouter) PUT(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.PutWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) PutWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)
	r.handle(http.MethodPut, path, h, routerOptions)
}

func (r *openapiRouter) OPTIONS(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.OptionsWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) OptionsWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)
	r.handle(http.MethodOptions, path, h, routerOptions)
}

func (r *openapiRouter) HEAD(path string, h handlerGenerator, options ...OpenAPIOption) {
	r.HeadWithOptions(path, h, WithOpenAPIOptions(options...))
}

func (r *openapiRouter) HeadWithOptions(path string, h handlerGenerator, options ...RouterOption) {
	routerOptions := mergeRouterOptions(options...)
	r.handle(http.MethodHead, path, h, routerOptions)
}
