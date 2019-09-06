package handler

import (
	"fmt"
	"github.com/ryotarai/github-api-auth-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type dummyAuthz struct {
	allowed bool
}

func (a dummyAuthz) IsRequestAllowed(username string, r *http.Request) (bool, error) {
	return a.allowed, nil
}

func dummyOrigin() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Authorization: %s\nMethod: %s\nPath: %s\n",
			r.Header.Get("Authorization"),
			r.Method,
			r.URL.Path)
	})
}

func TestServeHTTP(t *testing.T) {
	origin := httptest.NewServer(dummyOrigin())

	pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	authz := dummyAuthz{allowed: true}
	cfg := &config.Config{
		Passwords: map[string][]string{
			"user1": []string{
				string(pw),
			},
		},
	}
	originURL, err := url.Parse(origin.URL)
	assert.NoError(t, err)
	accessToken := "THISISACCESSTOKEN"

	h, err := New(cfg, originURL, accessToken, authz)
	assert.NoError(t, err)

	cases := []struct{
		password string
		code int
		body string
	}{
		{password: "password1", code: 200, body: "Authorization: token THISISACCESSTOKEN\nMethod: POST\nPath: /a/b/c\n"},
		{password: "invalid", code: 401, body: ""},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/a/b/c", strings.NewReader(""))
		r.SetBasicAuth("user1", c.password)
		h.ServeHTTP(w, r)
		out, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)
		assert.Equal(t, c.code, w.Code)
		assert.Equal(t, c.body, string(out))
	}
}
