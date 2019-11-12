package render_test

import (
	"context"
	"fmt"
	"github.com/fjgal/go-chi-jsonapi/render"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)
import chi_render "github.com/go-chi/render"

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name               string
		contentTypeString  string
		expectedContetType chi_render.ContentType
	}{
		{
			name:               "jsonapi",
			contentTypeString:  "application/vnd.api+json",
			expectedContetType: render.ContentTypeJSONAPI,
		},
		{
			name:               "defaults to go-chi/render (json)",
			contentTypeString:  "application/json",
			expectedContetType: chi_render.ContentTypeJSON,
		},
		{
			name:               "defaults to go-chi/render (unknown",
			contentTypeString:  "application/bogus",
			expectedContetType: chi_render.ContentTypeUnknown,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedContetType, render.GetContentType(test.contentTypeString))
		})
	}
}

func TestGetRequestContentType(t *testing.T) {
	tests := []struct {
		name               string
		contentTypeHeader  string
		contentTypeContext chi_render.ContentType
		expectedContetType chi_render.ContentType
	}{
		{
			name:               "jsonapi",
			contentTypeHeader:  "application/vnd.api+json",
			contentTypeContext: render.ContentTypeJSONAPI,
			expectedContetType: render.ContentTypeJSONAPI,
		},
		{
			name:               "defaults to go-chi/render (json)",
			contentTypeHeader:  "application/json",
			contentTypeContext: chi_render.ContentTypeJSON,
			expectedContetType: chi_render.ContentTypeJSON,
		},
		{
			name:               "defaults to go-chi/render (unknown)",
			contentTypeHeader:  "application/bogus",
			contentTypeContext: chi_render.ContentTypeUnknown,
			expectedContetType: chi_render.ContentTypeUnknown,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s - from Content-Type header", test.name), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			r.Header.Set("Content-Type", test.contentTypeHeader)
			assert.Equal(t, test.expectedContetType, render.GetRequestContentType(r))
		})
		t.Run(fmt.Sprintf("%s - from context", test.name), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			r = r.WithContext(context.WithValue(r.Context(), chi_render.ContentTypeCtxKey, test.contentTypeContext))
			assert.Equal(t, test.expectedContetType, render.GetRequestContentType(r))
		})
	}
}

type mockHandler struct {
	r http.Request
}

func (h mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", string(r.Context().Value(chi_render.ContentTypeCtxKey).(chi_render.ContentType)))
}

// TestSetContentType ensures that go-chi/render SetContentType middleware works with JSON API content type
func TestSetContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType chi_render.ContentType
	}{
		{
			name:        "jsonapi",
			contentType: render.ContentTypeJSONAPI,
		},
		{
			name:        "defaults to go-chi/render (json)",
			contentType: chi_render.ContentTypeJSON,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextHandler := mockHandler{}
			mw := chi_render.SetContentType(test.contentType)
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			w := httptest.NewRecorder()
			mw(nextHandler).ServeHTTP(w, r)
			assert.Equal(t, string(test.contentType), w.Header().Get("Content-Type"))
		})
	}
}

func TestGetAcceptedContentType(t *testing.T) {
	tests := []struct {
		name                string
		acceptHeader        string
		expectedContentType chi_render.ContentType
	}{
		{
			name:                "jsonapi",
			expectedContentType: render.ContentTypeJSONAPI,
			acceptHeader:        "application/vnd.api+json",
		},
		{
			name:                "defaults to go-chi/render (json)",
			acceptHeader:        "application/json",
			expectedContentType: chi_render.ContentTypeJSON,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			r.Header.Set("Accept", test.acceptHeader)
			assert.Equal(t, test.expectedContentType, render.GetAcceptedContentType(r))
		})
	}
}
