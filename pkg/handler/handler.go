package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/ryotarai/github-api-auth-proxy/pkg/authz"
	"github.com/ryotarai/github-api-auth-proxy/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	config      *config.Config
	originURL   *url.URL
	accessToken string
	authz       authz.Client
}

func New(config *config.Config, originURL *url.URL, accessToken string, authz authz.Client) (*Handler, error) {
	return &Handler{
		config:      config,
		originURL:   originURL,
		accessToken: accessToken,
		authz:       authz,
	}, nil
}

func (h *Handler) authn(username, password string) bool {
	hashedPasswords, ok := h.config.Passwords[username]
	if !ok {
		return false
	}

	for _, hashed := range hashedPasswords {
		err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		if err == nil {
			return true
		}
		log.Printf("Error from CompareHashAndPassword: %s", err)
	}
	return false
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := getCredFromRequest(r)
	if !ok {
		log.Printf("WARN: Failed to get both Basic Auth credential and Authorization token, url: %s\n", r.URL)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !h.authn(username, password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allowed, err := h.authz.IsRequestAllowed(username, r)
	if err != nil {
		log.Printf("ERR: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !allowed {
		log.Printf("WARN: request was not allowed, url: %s, username: %s\n", r.URL, username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
		// https://github.blog/2012-09-21-easier-builds-and-deployments-using-git-over-https-and-oauth/
		r.SetBasicAuth(h.accessToken, "x-oauth-basic")
	} else {
		r.Header.Set("Authorization", fmt.Sprintf("token %s", h.accessToken))
	}
	r.Host = h.originURL.Host

	httputil.NewSingleHostReverseProxy(h.originURL).ServeHTTP(w, r)
}

func getCredFromRequest(r *http.Request) (username string, password string, ok bool) {
	username, password, ok = r.BasicAuth()
	if ok {
		return
	}

	username, password, ok = getCredFromAuthorizationToken(r, "token")
	if ok {
		return
	}

	username, password, ok = getCredFromAuthorizationToken(r, "bearer")
	if ok {
		return
	}

	return
}

func getCredFromAuthorizationToken(r *http.Request, authType string) (username string, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}

	prefix := fmt.Sprintf("%s ", authType)
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	usernamePassword := strings.SplitN(auth[len(prefix):], ":", 2)
	if len(usernamePassword) < 2 {
		return
	}
	return usernamePassword[0], usernamePassword[1], true
}
