package cli

import (
	"flag"
	"github.com/ryotarai/github-api-auth-proxy/pkg/config"
	"github.com/ryotarai/github-api-auth-proxy/pkg/handler"
	"github.com/ryotarai/github-api-auth-proxy/pkg/authz"
	"log"
	"net/http"
	"net/url"
)

type CLI struct {
}

func New() *CLI {
	return &CLI{}
}

func (c *CLI) Start(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	listen := fs.String("listen", ":8080", "")
	originURLStr := fs.String("origin-url", "", "")
	opaServerURLStr := fs.String("authz-server-url", "", "")
	accessToken := fs.String("access-token", "", "")
	configPath := fs.String("config", "", "")

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	originURL, err := url.Parse(*originURLStr)
	if err != nil {
		return err
	}

	opaServerURL, err := url.Parse(*opaServerURLStr)
	if err != nil {
		return err
	}

	cfg, err := config.LoadYAMLFile(*configPath)
	if err != nil {
		return err
	}

	authz := authz.NewOPAClient(opaServerURL)

	h, err := handler.New(cfg, originURL, *accessToken, authz)
	if err != nil {
		return err
	}

	log.Printf("INFO: Listening on %s", *listen)
	err = http.ListenAndServe(*listen, h)
	if err != nil {
		return err
	}

	return nil
}
