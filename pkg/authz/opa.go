package authz

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"io/ioutil"
	"net/http"
)

type OPAClient struct {
	policyFile string
	policyBody string
}

func NewOPAClient(policyFile string) (*OPAClient, error) {
	bs, err := ioutil.ReadFile(policyFile)
	if err != nil {
		return nil, err
	}

	return &OPAClient{
		policyFile: policyFile,
		policyBody: string(bs),
	}, nil
}

func (c *OPAClient) IsRequestAllowed(username string, r *http.Request) (bool, error) {
	input := map[string]interface{}{
		"username": username,
		"method":   r.Method,
		"path":     r.URL.Path,
		"query":    r.URL.Query(),
		"header":   r.Header,
	}

	eval := rego.New(
		rego.Query("data.github.authz.allow"),
		rego.Input(input),
		rego.Module(c.policyFile, c.policyBody),
	)

	rs, err := eval.Eval(context.TODO())
	if err != nil {
		return false, err
	}

	if len(rs) == 0 {
		return false, fmt.Errorf("decision is undefined")
	}

	fmt.Printf("%+v\n", rs)
	allowed, ok := rs[0].Expressions[0].Value.(bool)
	if !ok {
		return false, fmt.Errorf("invalid policy")
	}

	return allowed, nil
}
