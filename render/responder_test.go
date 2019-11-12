package render_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/fjgal/go-chi-jsonapi/render"
	chi_render "github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONAPI(t *testing.T) {

	t.Run("should render a proper json api payload", func(t *testing.T) {

	})
}

func TestDefaultResponder(t *testing.T) {

	tests := []struct {
		name               string
		contentTypeContext chi_render.ContentType
		contentTypeHeader  string
		expectedBody       []byte
	}{
		{
			name:               "json api",
			contentTypeContext: render.ContentTypeJSONAPI,
			contentTypeHeader:  "application/vnd.api+json",
			expectedBody:       []byte(`{"data":{"type":"blogs","id":"11","attributes":{"current_post_id":0,"title":"The Best Blog","view_count":0},"relationships":{"current_post":{"data":null},"posts":{"data":[]}}}}`),
		},
		{
			name:               "defaults to go-chi/render",
			contentTypeContext: chi_render.ContentTypeJSON,
			contentTypeHeader:  "application/json",
			expectedBody:       []byte(`{"ID":11,"Title":"The Best Blog","Posts":null,"CurrentPost":null,"CurrentPostID":0,"CreatedAt":"0001-01-01T00:00:00Z","ViewCount":0}`),
		},
	}

	v := Blog{
		ID:    11,
		Title: "The Best Blog",
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s - with Accept header", test.name), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			r.Header.Set("Accept", test.contentTypeHeader)
			w := httptest.NewRecorder()
			render.DefaultResponder(w, r, &v)
			assert.Equal(t, string(test.expectedBody), strings.TrimSpace(w.Body.String()))
		})
		t.Run(fmt.Sprintf("%s - with content type in request context", test.name), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			r = r.WithContext(context.WithValue(r.Context(), chi_render.ContentTypeCtxKey, test.contentTypeContext))
			w := httptest.NewRecorder()
			render.DefaultResponder(w, r, &v)
			assert.Equal(t, string(test.expectedBody), strings.TrimSpace(w.Body.String()))
		})
	}

}

func TestDefaultResponder_Errors(t *testing.T) {

	tests := []struct {
		name         string
		contentType  chi_render.ContentType
		status       int
		err          error
		expectedBody []byte
	}{
		{
			name:         "json api with error payload - should render as error object (422)",
			contentType:  render.ContentTypeJSONAPI,
			status:       http.StatusUnprocessableEntity,
			err:          errors.New("something went wrong"),
			expectedBody: []byte(`{"errors":[{"title":"Unprocessable Entity","detail":"something went wrong","status":"422"}]}`),
		},
		{
			name:         "json api with error payload - should render as error object (500)",
			contentType:  render.ContentTypeJSONAPI,
			status:       http.StatusInternalServerError,
			err:          errors.New("something went wrong"),
			expectedBody: []byte(`{"errors":[{"title":"Internal Server Error","detail":"something went wrong","status":"500"}]}`),
		},
		{
			name:         "json defaults to go-chi/render - doesn't render errors in any particular way",
			contentType:  chi_render.ContentTypeJSON,
			status:       500,
			err:          errors.New("something went wrong"),
			expectedBody: []byte(`{}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
			r = r.WithContext(context.WithValue(r.Context(), chi_render.ContentTypeCtxKey, test.contentType))
			chi_render.Status(r, test.status)
			w := httptest.NewRecorder()
			render.DefaultResponder(w, r, test.err)
			assert.Equal(t, string(test.expectedBody), strings.TrimSpace(w.Body.String()))
			assert.Equal(t, test.status, w.Result().StatusCode)
		})
	}

}
