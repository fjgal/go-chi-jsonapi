package render

import (
	chi_render "github.com/go-chi/render"
	"github.com/google/jsonapi"
	"io"
	"io/ioutil"
	"net/http"
)

// DefaultDecoder decodes JSON API
// delegates any other content type to github.con/go-chi/render
func DefaultDecoder(r *http.Request, v interface{}) error {
	var err error

	switch GetRequestContentType(r) {
	case ContentTypeJSONAPI:
		err = DecodeJSONAPI(r.Body, v)
	default:
		err = chi_render.DefaultDecoder(r, v)
	}

	return err
}

func DecodeJSONAPI(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return jsonapi.UnmarshalPayload(r, v)
}
