package server

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func EncodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if isGZip := strings.Contains(strings.ToLower(w.Header().Get("Accept-Encoding")), "gzip"); isGZip {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		if b, ok := response.([]byte); ok && len(b) > 0 {
			_, err := gz.Write(b)
			return err
		}
		return json.NewEncoder(gz).Encode(response)
	} else if b, ok := response.([]byte); ok && len(b) > 0 {
		_, err := w.Write(b)
		return err
	}
	return json.NewEncoder(w).Encode(response)
}
