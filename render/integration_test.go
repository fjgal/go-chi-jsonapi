package render_test

import (
	"bytes"
	"fmt"
	jsonapi_render "github.com/fjgal/go-chi-jsonapi/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	tests := []struct {
		name              string
		contentTypeHeader string
		acceptHeader      string
		reqBody           []byte
		respBody          []byte
		status            int
	}{
		{
			name:              "jsonapi in / jsonapi out",
			contentTypeHeader: "application/vnd.api+json",
			acceptHeader:      "application/vnd.api+json",
			reqBody:           []byte(`{"data":{"type":"blogs","id":"42","attributes":{"current_post_id":0,"title":"The Best Blog","view_count":0},"relationships":{"current_post":{"data":null},"posts":{"data":[]}}}}`),
			respBody:          []byte(`{"data":{"type":"blogs","id":"42","attributes":{"current_post_id":0,"title":"The Best Blog","view_count":0},"relationships":{"current_post":{"data":null},"posts":{"data":[]}}}}`),
			status:            202,
		},
		{
			name:              "jsonapi in / json out",
			contentTypeHeader: "application/vnd.api+json",
			acceptHeader:      "application/json",
			reqBody:           []byte(`{"data":{"type":"blogs","id":"42","attributes":{"current_post_id":0,"title":"The Best Blog","view_count":0},"relationships":{"current_post":{"data":null},"posts":{"data":[]}}}}`),
			respBody:          []byte(`{"ID":42,"Title":"The Best Blog","Posts":null,"CurrentPost":null,"CurrentPostID":0,"CreatedAt":"0001-01-01T00:00:00Z","ViewCount":0}`),
			status:            202,
		},
		{
			name:              "json in / jsonapi out",
			contentTypeHeader: "application/json",
			acceptHeader:      "application/vnd.api+json",
			reqBody:           []byte(`{"ID":42,"Title":"The Best Blog","Posts":null,"CurrentPost":null,"CurrentPostID":0,"CreatedAt":"0001-01-01T00:00:00Z","ViewCount":0}`),
			respBody:          []byte(`{"data":{"type":"blogs","id":"42","attributes":{"current_post_id":0,"title":"The Best Blog","view_count":0},"relationships":{"current_post":{"data":null},"posts":{"data":[]}}}}`),
			status:            202,
		},
		{
			name:              "non-decodable request body",
			contentTypeHeader: "application/json",
			acceptHeader:      "application/vnd.api+json",
			reqBody:           []byte(`{"data":{{{`),
			respBody:          []byte(`{"errors":[{"title":"Unprocessable Entity","detail":"invalid character '{' looking for beginning of object key string","status":"422"}]}`),
			status:            422, // unprocessable entity
		},
	}

	// plug jsonapi responder and decoder
	render.Respond = jsonapi_render.DefaultResponder
	render.Decode = jsonapi_render.DefaultDecoder

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s - decode/respond", test.name), func(t *testing.T) {
			router := chi.NewRouter()

			// handler that decodes and echoes back the same object
			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				var blog Blog
				err := render.Decode(r, &blog)
				if err != nil {
					t.Log(err)
					render.Status(r, http.StatusUnprocessableEntity)
					render.Respond(w, r, err)
					return
				}
				render.Status(r, http.StatusAccepted)
				render.Respond(w, r, &blog)
			})

			ts := httptest.NewServer(router)

			r, err := http.NewRequest(http.MethodGet, ts.URL, bytes.NewReader(test.reqBody))
			if err != nil {
				t.Log(err)
				t.FailNow()
			}

			r.Header.Set("Content-Type", test.contentTypeHeader)
			r.Header.Set("Accept", test.acceptHeader)

			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}

			assert.Equal(t, test.status, resp.StatusCode)
			b, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, string(test.respBody), strings.TrimSpace(string(b)))
		})

		t.Run(fmt.Sprintf("%s - bind/render", test.name), func(t *testing.T) {
			router := chi.NewRouter()

			// handler that decodes and echoes back the same object
			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				var blog Blog
				err := render.Bind(r, &blog)
				if err != nil {
					t.Log(err)
					render.Status(r, http.StatusUnprocessableEntity)
					render.Respond(w, r, err)
					return
				}
				render.Status(r, http.StatusAccepted)
				err = render.Render(w, r, &blog)
				if err != nil {
					t.Log(err)
					return
				}
			})

			ts := httptest.NewServer(router)

			r, err := http.NewRequest(http.MethodGet, ts.URL, bytes.NewReader(test.reqBody))
			if err != nil {
				t.Log(err)
				t.FailNow()
			}

			r.Header.Set("Content-Type", test.contentTypeHeader)
			r.Header.Set("Accept", test.acceptHeader)

			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}

			assert.Equal(t, test.status, resp.StatusCode)
			b, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, string(test.respBody), strings.TrimSpace(string(b)))
		})
	}

}
