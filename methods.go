package httpy

import "net/http"

// Get returns a new Request used to make a GET request to the given url.
func Get(url string) *Request {
	return NewRequest(url, http.MethodGet)
}

// Post returns a new Request used to make POST requests to the given url.
func Post(url string) *Request {
	return NewRequest(url, http.MethodPost)
}

// Put returns a new Request used to make PUT requests to the given url.
func Put(url string) *Request {
	return NewRequest(url, http.MethodPut)
}

// Patch returns a new Request used to make PATCH requests to the given url.
func Patch(url string) *Request {
	return NewRequest(url, http.MethodPatch)
}

// Delete returns a new Request used to make DELETE requests to the given url.
func Delete(url string) *Request {
	return NewRequest(url, http.MethodDelete)
}
