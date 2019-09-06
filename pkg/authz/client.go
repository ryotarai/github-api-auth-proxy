package authz

import "net/http"

type Client interface {
	IsRequestAllowed(username string, r *http.Request) (bool, error)
}
