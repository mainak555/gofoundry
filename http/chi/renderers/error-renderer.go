package renderers

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	ReasonPhrase   string `json:"reason"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
	ErrorObject    any    `json:"details,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		ReasonPhrase:   "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		ReasonPhrase:   "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrUnAuthorized = &ErrResponse{HTTPStatusCode: 401, ReasonPhrase: "unauthorize"}
var ErrBadRequest = &ErrResponse{HTTPStatusCode: 400, ReasonPhrase: "Invalid request!"}
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, ReasonPhrase: "Resource not found!"}
var ErrMissingHeader = &ErrResponse{HTTPStatusCode: 400, ReasonPhrase: "Header Part Missing!"}
