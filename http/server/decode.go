package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"gofoundry/dtos"
)

func DecodeRequestBody[T interface{}](_ context.Context, r *http.Request) (interface{}, error) {
	var tModel T
	err := json.NewDecoder(r.Body).Decode(&tModel)
	return tModel, err
}

func DecodeRequestAsJsonString(_ context.Context, r *http.Request) (interface{}, error) {
	byteBody, err := io.ReadAll(r.Body)
	return string(byteBody), err
}

func PropageRawRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return *r, nil
}

func DecodeRequestNone(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeRequestQuery(ctx context.Context, r *http.Request) (interface{}, error) {
	return DecodeRequestBody[dtos.QueryRequestDTO](ctx, r)
}
