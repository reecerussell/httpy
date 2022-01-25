package httpy

import (
	"encoding/json"
	"net/http"
)

// Response represents the response from a HTTP request. It wraps http.Response,
// while providing helpful functions alongside.
type Response http.Response

// DecodeJSON the Response's JSON body to the given dest.
func (r *Response) DecodeJSON(dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}
