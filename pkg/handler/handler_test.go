package handler

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ryotarai/github-api-auth-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type dummyAuthz struct {
	allowed bool
}

func (a dummyAuthz) IsRequestAllowed(username string, r *http.Request) (bool, error) {
	return a.allowed, nil
}

func dummyOrigin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	cases := []struct {
		reqFunc  func() *http.Request
		password string
		code     int
		body     string
	}{
		// Token
		{
			reqFunc: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/v3/user", strings.NewReader(""))
				r.Header.Set("Authorization", fmt.Sprintf("token %s", "user1:password1"))
				return r
			},
			code: 200,
			body: fmt.Sprintf("Authorization: token " + accessToken + "\nMethod: GET\nPath: /api/v3/user\n"),
		},
		{
			reqFunc: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/v3/user", strings.NewReader(""))
				r.Header.Set("Authorization", fmt.Sprintf("bearer %s", "user1:password1"))
				return r
			},
			code: 200,
			body: fmt.Sprintf("Authorization: token " + accessToken + "\nMethod: GET\nPath: /api/v3/user\n"),
		},
		{
			reqFunc: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/v3/user", strings.NewReader(""))
				r.Header.Set("Authorization", fmt.Sprintf("bearer %s", "invalid"))
				return r
			},
			code: 401,
			body: "",
		},
		// Basic Authorization
		{
			reqFunc: func() *http.Request {
				r := httptest.NewRequest("POST", "/a/b/c", strings.NewReader(""))
				r.SetBasicAuth("user1", "password1")
				return r
			},
			code: 200,
			body: fmt.Sprintf("Authorization: Basic %s\nMethod: POST\nPath: /a/b/c\n", base64.StdEncoding.EncodeToString([]byte(accessToken+":x-oauth-basic"))),
		},
		{
			reqFunc: func() *http.Request {
				r := httptest.NewRequest("POST", "/a/b/c", strings.NewReader(""))
				r.SetBasicAuth("user1", "invalid")
				return r
			},
			code: 401,
			body: "",
		},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		r := c.reqFunc()
		h.ServeHTTP(w, r)
		out, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)
		assert.Equal(t, c.code, w.Code)
		assert.Equal(t, c.body, string(out))
	}
}
