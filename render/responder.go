package render

import (
	"bytes"
	chi_render "github.com/go-chi/render"
	"github.com/google/jsonapi"
	"net/http"
	"strconv"
)

// Respond handles JSON API responses and delegates any other content type to github.com/go-chi/render
// automatically setting the Content-Type based on request headers
func DefaultResponder(w http.ResponseWriter, r *http.Request, v interface{}) {

	// Format response based on request Accept header.
	switch GetAcceptedContentType(r) {
	case ContentTypeJSONAPI:
		JSONAPI(w, r, v)
	default:
		chi_render.DefaultResponder(w, r, v)
	}
}

// JSONAPI marshals `v` to JSONAPI, automatically setting Content-Type as application/vnd.api+json
func JSONAPI(w http.ResponseWriter, r *http.Request, v interface{}) {

	w.Header().Set("Content-Type", "application/vnd.api+json")

	switch v.(type) {
	case error:
		renderError(w, r, v.(error))
	default:
		renderPayload(w, r, v)
	}

}

func toJSONAPIErrors(status int, errors ...error) (jsonapierrors []*jsonapi.ErrorObject) {
	for _, e := range errors {
		jsonapierrors = append(jsonapierrors, &jsonapi.ErrorObject{
			Title:  http.StatusText(status),
			Detail: e.Error(),
			Status: strconv.Itoa(status),
			Code:   "",
			Meta:   nil,
		})
	}
	return
}

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError // default if not set in context
	if i, ok := r.Context().Value(chi_render.StatusCtxKey).(int); ok {
		status = i
	}
	w.WriteHeader(status)
	_ = jsonapi.MarshalErrors(w, toJSONAPIErrors(status, err))
}

func renderPayload(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	if err := jsonapi.MarshalPayload(buf, v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = jsonapi.MarshalErrors(w, toJSONAPIErrors(http.StatusInternalServerError, err))
		return
	}

	if status, ok := r.Context().Value(chi_render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	_, _ = w.Write(buf.Bytes())

}
