package authz

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOPAClient(t *testing.T) {
	var allowed bool
	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if allowed {
			fmt.Fprint(w, `{"result": {"allow": true}}\n`)
		} else {
			fmt.Fprint(w, `{"result": {"allow": false}}\n`)
		}
	}))

	serverURL, err := url.Parse(server.URL)
	assert.NoError(t, err)

	c := NewOPAClient(serverURL)

	allowed = true
	returned, err := c.IsRequestAllowed("user1", httptest.NewRequest(
		"POST", "/a/b/c", strings.NewReader(`{"a": "b"}`)))
	assert.NoError(t, err)
	assert.Equal(t, true, returned)

	allowed = false
	returned, err = c.IsRequestAllowed("user1", httptest.NewRequest(
		"POST", "/a/b/c", strings.NewReader(`{"a": "b"}`)))
	assert.NoError(t, err)
	assert.Equal(t, false, returned)
}