package generics

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
)

type IEndpoints interface {
	GetAllEndpoint() endpoint.Endpoint
	DeleteEndpoint() endpoint.Endpoint
	CreateOneEndpoint() endpoint.Endpoint
	GetByIdEndpoint() endpoint.Endpoint
	UpdateByIdEndpoint() endpoint.Endpoint
	DeleteByIdEndpoint() endpoint.Endpoint
}

type IApiRouterRequestDecoder interface {
	GetAllDeoder() http.DecodeRequestFunc
	GetByIdDecoder() http.DecodeRequestFunc
	DeleteByIdDecoder() http.DecodeRequestFunc
	DeleteManyDecoder() http.DecodeRequestFunc
}
