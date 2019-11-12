package render_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fjgal/go-chi-jsonapi/render"
	chi_render "github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultDecoder(t *testing.T) {
	tests := []struct {
		name               string
		contentTypeContext chi_render.ContentType
		contentTypeHeader  string
		body               []byte
	}{
		{
			name:               "json api",
			contentTypeContext: render.ContentTypeJSONAPI,
			contentTypeHeader:  "application/vnd.api+json",
			body:               []byte(`{"data":{"type":"blogs","id":"11","attributes":{"current_post_id":0,"title":"The Best Blog","view_count":0},"relationships":{"current_post":{"data":null},"posts":{"data":[]}}}}`),
		},
		{
			name:               "defaults to go-chi/render",
			contentTypeContext: chi_render.ContentTypeJSON,
			contentTypeHeader:  "application/json",
			body:               []byte(`{"ID":11,"Title":"The Best Blog","Posts":null,"CurrentPost":null,"CurrentPostID":0,"CreatedAt":"0001-01-01T00:00:00Z","ViewCount":0}`),
		},
	}

	expectedValue := Blog{
		ID:    11,
		Title: "The Best Blog",
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s - with Content-Type header", test.name), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", bytes.NewReader(test.body))
			r.Header.Set("Content-Type", test.contentTypeHeader)
			var v Blog
			err := render.DefaultDecoder(r, &v)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}
			assert.Equal(t, expectedValue, v)
		})
		t.Run(fmt.Sprintf("%s - with content type in request context", test.name), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://www.example.com", bytes.NewReader(test.body))
			r = r.WithContext(context.WithValue(r.Context(), chi_render.ContentTypeCtxKey, test.contentTypeContext))
			var v Blog
			err := render.DefaultDecoder(r, &v)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}
			assert.Equal(t, expectedValue, v)
		})
	}
}
