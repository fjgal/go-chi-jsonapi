package render

import (
	"net/http"
	"strings"
)
import chi_render "github.com/go-chi/render"

const (
	ContentTypeJSONAPI chi_render.ContentType = iota + 1000
)

// GetContentType extends go-ci/render to support application/vnd.api+json
func GetContentType(s string) chi_render.ContentType {
	s = strings.TrimSpace(strings.Split(s, ";")[0])
	switch s {
	case "application/vnd.api+json":
		return ContentTypeJSONAPI
	default:
		return chi_render.GetContentType(s)
	}
}

func GetAcceptedContentType(r *http.Request) chi_render.ContentType {
	if contentType, ok := r.Context().Value(chi_render.ContentTypeCtxKey).(chi_render.ContentType); ok {
		return contentType
	}

	var contentType chi_render.ContentType

	// Parse request Accept header.
	fields := strings.Split(r.Header.Get("Accept"), ",")
	if len(fields) > 0 {
		contentType = GetContentType(strings.TrimSpace(fields[0]))
	}

	if contentType == chi_render.ContentTypeUnknown {
		contentType = chi_render.ContentTypePlainText
	}

	return contentType
}

// GetRequestContentType is a helper function that returns ContentType based on
// context or request headers.
func GetRequestContentType(r *http.Request) chi_render.ContentType {
	if contentType, ok := r.Context().Value(chi_render.ContentTypeCtxKey).(chi_render.ContentType); ok {
		return contentType
	}
	return GetContentType(r.Header.Get("Content-Type"))
}
