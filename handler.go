package httpserver

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ihezebin/openapi"
	"github.com/pkg/errors"

	"github.com/ihezebin/soup/logger"
)

type Handler[RequestT any, ResponseT any] func(c *gin.Context, req RequestT) (resp ResponseT, err error)

func newErrHandlerFunc[RequestT any, ResponseT any](err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data ResponseT
		body := &Body[ResponseT]{
			Data: data,
		}
		errx := ErrorWithInternalServer()
		errx.Err = err

		body.WithErr(errx)
		c.AbortWithStatusJSON(body.status, body)
	}
}

func newGinHandlerFunc[RequestT any, ResponseT any](handler Handler[RequestT, ResponseT], isRequestStruct bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var err error

		body := &Body[ResponseT]{
			status: http.StatusOK,
			Code:   CodeOK,
		}

		// requestPtr := reflect.New(reflect.TypeOf((*RequestT)(nil)).Elem()).Interface()
		requestPtr := new(RequestT)
		if isRequestStruct { // 如果是结构体可自由绑定
			if err = c.ShouldBind(requestPtr); err != nil {
				logger.WithError(err).Errorf(ctx, "failed to bind, uri: %s", c.Request.RequestURI)
				body = body.WithErr(ErrorWithInternalServer())
				c.PureJSON(body.status, body)
				return
			}
		} else { // 如果是 map 则优先从 body 读, 如果 body 为空则从 query 读
			if err = c.ShouldBindWith(requestPtr, mapBinding{}); err != nil {
				logger.WithError(err).Errorf(ctx, "failed to bind, uri: %s", c.Request.RequestURI)
				body = body.WithErr(ErrorWithBadRequest())
				c.PureJSON(http.StatusBadRequest, body)
				return
			}
		}

		if isRequestStruct && len(c.Params) > 0 {
			if err = c.ShouldBindUri(requestPtr); err != nil {
				logger.WithError(err).Errorf(ctx, "failed to bind uri, uri: %s, param: %+v", c.Request.RequestURI, c.Params)
				body = body.WithErr(ErrorWithBadRequest())
				c.PureJSON(http.StatusBadRequest, body)
				return
			}
		}

		if isRequestStruct && len(c.Request.Header) > 0 {
			if err = c.ShouldBindHeader(requestPtr); err != nil {
				logger.WithError(err).Errorf(ctx, "failed to bind header, uri: %s, header: %+v", c.Request.RequestURI, c.Request.Header)
				body = body.WithErr(ErrorWithBadRequest())
				c.PureJSON(http.StatusBadRequest, body)
				return
			}
		}

		var response ResponseT
		response, err = handler(c, *requestPtr)
		if c.Writer.Written() {
			return
		}

		// handle error
		if err != nil {
			var errx *Err
			if errors.As(err, &errx) {
				body = body.WithErr(errx)
			} else {
				body = body.WithErr(ErrorWithInternalServer())
			}
		} else {
			body.Data = response
		}

		// handle success
		c.PureJSON(body.status, body)
	}
}

// handlerGenerator return [requestBody, responseBody, query, params, requestHeader, responseHeader, gin.HandlerFunc]
type handlerGenerator func() (*openapi.Model, *openapi.Model, map[string]openapi.QueryParam, map[string]openapi.PathParam, map[string]openapi.HeaderParam, map[string]openapi.HeaderParam, gin.HandlerFunc)

func NewHandler[RequestT any, ResponseT any](handler Handler[RequestT, ResponseT]) handlerGenerator {
	return func() (*openapi.Model, *openapi.Model, map[string]openapi.QueryParam, map[string]openapi.PathParam, map[string]openapi.HeaderParam, map[string]openapi.HeaderParam, gin.HandlerFunc) {
		responseBodyModel := openapi.ModelOf[Body[ResponseT]]()

		request := new(RequestT)
		requestType := reflect.TypeOf(request).Elem()

		for requestType.Kind() == reflect.Ptr {
			requestType = requestType.Elem()
		}

		//  RequestT 必须是结构体或者 map
		if requestType.Kind() != reflect.Struct && requestType.Kind() != reflect.Map {
			err := fmt.Errorf("request type must be struct or map, but got %T", request)
			return nil, &responseBodyModel, nil, nil, nil, nil, newErrHandlerFunc[RequestT, ResponseT](err)
		}

		isRequestStruct := requestType.Kind() == reflect.Struct
		ginHandleFunc := newGinHandlerFunc(handler, isRequestStruct)
		if requestType.Kind() == reflect.Map {
			requestBodyModel := openapi.ModelFromType(requestType)
			return &requestBodyModel, &responseBodyModel, nil, nil, nil, nil, ginHandleFunc
		}

		// 通过反射获取 request 的字段和 tag
		var requestBodyStructFields []reflect.StructField
		query := map[string]openapi.QueryParam{}
		params := map[string]openapi.PathParam{}
		requestHeader := map[string]openapi.HeaderParam{}
		responseHeader := map[string]openapi.HeaderParam{}

		for i := 0; i < requestType.NumField(); i++ {
			field := requestType.Field(i)
			fieldType := field.Type

			tagDescription := field.Tag.Get("description")
			if tagDescription == "" {
				tagDescription = field.Tag.Get("desc")
			}

			isRequired := false
			isAllowEmpty := false
			tagOpenApi := field.Tag.Get("openapi")
			if tagOpenApi != "" { // required,empty
				parts := strings.Split(tagOpenApi, ",")
				for _, part := range parts {
					if part == "required" {
						isRequired = true
					}
					if part == "empty" {
						isAllowEmpty = true
					}
				}
			}

			tagJson := field.Tag.Get("json")
			if tagJson != "" {
				requestBodyStructFields = append(requestBodyStructFields, field)
			}

			primitiveType := openapi.PrimitiveTypeString
			switch fieldType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				primitiveType = openapi.PrimitiveTypeInteger
			case reflect.Bool:
				primitiveType = openapi.PrimitiveTypeBool
			case reflect.Float64, reflect.Float32:
				primitiveType = openapi.PrimitiveTypeFloat64
			}

			tagQuery := field.Tag.Get("query")
			if tagQuery == "" {
				tagQuery = field.Tag.Get("form")
			}
			if tagQuery != "" {
				queryParam := openapi.QueryParam{
					Description: tagDescription,
					Required:    isRequired,
					AllowEmpty:  isAllowEmpty,
					Type:        primitiveType,
				}

				query[tagQuery] = queryParam
			}

			tagUri := field.Tag.Get("uri")
			if tagUri != "" {
				param := openapi.PathParam{
					Description: tagDescription,
					Type:        primitiveType,
				}

				params[tagUri] = param
			}

			tagHeader := field.Tag.Get("header")
			if tagHeader != "" {
				headerParam := openapi.HeaderParam{
					Description: tagDescription,
					Type:        primitiveType,
					Required:    isRequired,
				}

				requestHeader[tagHeader] = headerParam
			}
		}
		var requestBodyModel *openapi.Model = nil

		if len(requestBodyStructFields) > 0 {
			requestBodyType := reflect.StructOf(requestBodyStructFields)
			a := requestBodyType.Kind()
			_ = a
			requestBodyModel = &openapi.Model{
				Type: requestBodyType,
			}
		}

		return requestBodyModel, &responseBodyModel, query, params, requestHeader, responseHeader, ginHandleFunc
	}
}
