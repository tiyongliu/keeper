package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Method int32

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

type ResponseStatistics struct {
	TraceId       string
	Msg           string
	ErrorType     string
	Step          string
	Cid           string
	RoutingList   []*BriefRouting
	DepartureCode string
	ArrivalCode   string
}

type BriefRouting struct {
	ProviderName string
	Carrier      string
}

type ApiResponse interface {
	ResponseStatus() string
	ErrorMsg() string
}

type PongResponse string

const DefaultPong = PongResponse("pong")

func (rs PongResponse) ResponseStatus() string {
	return ""
}

func (rs PongResponse) ErrorMsg() string {
	return ""
}

var InternalErrorRes = &ErrorResponse{Msg: "Internal Error", Status: 500}
var BadRequestErrorRes = &ErrorResponse{Msg: "Bad request", Status: 400}

type ErrorResponse struct {
	Status int
	Msg    string
}

func (rs *ErrorResponse) ResponseStatus() string {
	return strconv.Itoa(rs.Status)
}

func (rs *ErrorResponse) ErrorMsg() string {
	return rs.Msg
}

type DefaultSuccessResponse struct {
}

func (rs *DefaultSuccessResponse) ResponseStatus() string {
	return "200"
}

func (rs *DefaultSuccessResponse) ErrorMsg() string {
	return "Success"
}

type ApiHandlerFunc func(*gin.Context) (ApiResponse, error)

type ApiReqHandleAspect func(*gin.Context) error

type ApiResHandleAspect func(ctx *gin.Context, res interface{}) error

type ApiHandler struct {
	Method       Method
	Pattern      string
	BeforeHandle []ApiReqHandleAspect
	AfterHandle  []ApiResHandleAspect
	Handler      ApiHandlerFunc
}
type ApiDefinitionOption func(*ApiDefinition)

type ApiDefinition struct {
	sync.Mutex
	Method     Method
	Name       string
	PathPrefix string
	Handlers   map[string]*ApiHandler
}

func (method Method) String() string {
	switch method {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case DELETE:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func NewApiDefinition(options ...ApiDefinitionOption) *ApiDefinition {
	ad := &ApiDefinition{
		Handlers: map[string]*ApiHandler{},
	}

	for _, option := range options {
		option(ad)
	}

	return ad
}

func (ad *ApiDefinition) WithAssignHandle(method Method, pattern string, reqHandle []ApiReqHandleAspect, resHandle []ApiResHandleAspect, apiHandlerFunc ApiHandlerFunc) *ApiDefinition {
	ad.Lock()
	defer ad.Unlock()

	p := strings.ToLower(pattern)

	if strings.HasPrefix(p, "/") {
		p = p[1:]
	}

	_, ok := ad.Handlers[p]

	if ok {
		//logger.Fatal(logger.Message("Duplicate path[%s] for ApiDefinition[%s]", p, ad.Name))Name
	}
	ad.Handlers[p] = &ApiHandler{
		Method:       method,
		Pattern:      pattern,
		BeforeHandle: reqHandle,
		Handler:      apiHandlerFunc,
		AfterHandle:  resHandle,
	}

	return ad
}

func NewStandardReqHandle() []ApiReqHandleAspect {
	return []ApiReqHandleAspect{}
}

func NewStandardResHandle() []ApiResHandleAspect {
	return []ApiResHandleAspect{
		func(ctx *gin.Context, res interface{}) error {
			ctx.JSON(http.StatusOK, res)
			return nil
		},
	}
}

func WithName(name string) ApiDefinitionOption {
	return func(ad *ApiDefinition) {
		ad.Name = strings.ToLower(name)
	}
}

func WithPrefix(prefix string) ApiDefinitionOption {
	return func(ad *ApiDefinition) {
		ad.PathPrefix = strings.ToLower(prefix)
	}
}

func (ad *ApiDefinition) WithHandler(method Method, pattern string, apiHandlerFunc ApiHandlerFunc) *ApiDefinition {
	return ad.WithAssignHandle(method, pattern, NewStandardReqHandle(), NewStandardResHandle(), apiHandlerFunc)

}

func Ping(c *gin.Context) (ApiResponse, error) {
	return DefaultPong, nil
}
