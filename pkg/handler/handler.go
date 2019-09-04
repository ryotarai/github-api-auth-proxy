package handler

import (
	"fmt"
	"github.com/ryotarai/github-api-authz-proxy/pkg/opa"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Handler struct {
	originURL   *url.URL
	accessToken string
	authz       *opa.Client
}

func New(originURL *url.URL, accessToken string, authz *opa.Client) (*Handler, error) {
	return &Handler{
		originURL:   originURL,
		accessToken: accessToken,
		authz:       authz,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allowed, err := h.authz.IsRequestAllowed(r)
	if err != nil {
		log.Printf("ERR: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !allowed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Header.Set("Authorization", fmt.Sprintf("token %s", h.accessToken))
	r.Host = h.originURL.Host

	httputil.NewSingleHostReverseProxy(h.originURL).ServeHTTP(w, r)

	//r.URL.Scheme = h.originURL.Scheme
	//r.URL.Host = h.originURL.Host
	//r.URL.Path = path.Join(h.originURL.Path, r.URL.Path)
	//r.RequestURI = ""
	//
	//resp, err := http.DefaultClient.Do(r)
	//if err != nil {
	//	log.Printf("ERR: %s", err)
	//	w.WriteHeader(http.StatusBadGateway)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//w.WriteHeader(resp.StatusCode)
	//
	//header := w.Header()
	//for k, vs := range resp.Header {
	//	for _, v := range vs {
	//		header.Add(k, v)
	//	}
	//}
	//
	//_, err = io.Copy(w, resp.Body)
	//if err != nil {
	//	log.Printf("ERR: %s", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
}
