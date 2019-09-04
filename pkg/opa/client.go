package opa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	serverURL *url.URL
}

func NewClient(serverURL *url.URL) *Client {
	return &Client{
		serverURL: serverURL,
	}
}

type input struct {
	Input inputInput `json:"input"`
}

type inputInput struct {
	Username string      `json:"username"`
	Method   string      `json:"method"`
	Path     string      `json:"path"`
	Query    url.Values  `json:"query"`
	Header   http.Header `json:"header"`
	Body     interface{} `json:"body"`
}

type output struct {
	Result outputResult `json:"result"`
}

type outputResult struct {
	Allow bool `json:"allow"`
}

func (c *Client) IsRequestAllowed(username string, r *http.Request) (bool, error) {
	input := input{
		Input: inputInput{
			Username: username,
			Method:   r.Method,
			Path:     r.URL.Path,
			Query:    r.URL.Query(),
			Header:   r.Header,
		},
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

	if len(bodyBytes) > 0 {
		body := map[string]interface{}{}
		err = json.Unmarshal(bodyBytes, &body)
		if err != nil {
			return false, err
		}

		input.Input.Body = body
	}

	log.Printf("DEBUG: Input to OPA: %+v", input)

	u := *c.serverURL
	u.Path = path.Join(u.Path, "v1/data/httpapi/authz")

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return false, err
	}
	body := bytes.NewReader(inputJSON)

	resp, err := http.Post(u.String(), "application/json", body)
	if err != nil {
		return false, err
	}

	out := output{}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return false, err
	}
	log.Printf("DEBUG: Output from OPA: %+v", out)

	return out.Result.Allow, nil
}
