package httpy

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	r := Get("https://foo.bar")
	assert.Equal(t, http.MethodGet, r.method)
	assert.Equal(t, "https://foo.bar", r.url)
}

func TestPost(t *testing.T) {
	r := Post("https://foo.bar")
	assert.Equal(t, http.MethodPost, r.method)
	assert.Equal(t, "https://foo.bar", r.url)
}

func TestPut(t *testing.T) {
	r := Put("https://foo.bar")
	assert.Equal(t, http.MethodPut, r.method)
	assert.Equal(t, "https://foo.bar", r.url)
}

func TestPatch(t *testing.T) {
	r := Patch("https://foo.bar")
	assert.Equal(t, http.MethodPatch, r.method)
	assert.Equal(t, "https://foo.bar", r.url)
}

func TestDelete(t *testing.T) {
	r := Delete("https://foo.bar")
	assert.Equal(t, http.MethodDelete, r.method)
	assert.Equal(t, "https://foo.bar", r.url)
}
