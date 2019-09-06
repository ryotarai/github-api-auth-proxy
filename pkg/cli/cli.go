package cli

import (
	"flag"
	"fmt"
	"github.com/ryotarai/github-api-auth-proxy/pkg/authz"
	"github.com/ryotarai/github-api-auth-proxy/pkg/config"
	"github.com/ryotarai/github-api-auth-proxy/pkg/handler"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
	"syscall"
)

type CLI struct {
}

func New() *CLI {
	return &CLI{}
}

func (c *CLI) Start(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	listen := fs.String("listen", ":8080", "")
	configPath := fs.String("config", "", "")
	bcryptMode := fs.Bool("bcrypt", false, "")

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if *bcryptMode {
		generateBcrypt()
		return nil
	}

	cfg, err := config.LoadYAMLFile(*configPath)
	if err != nil {
		return err
	}

	originURL, err := url.Parse(cfg.OriginURL)
	if err != nil {
		return err
	}

	opaServerURL, err := url.Parse(cfg.OPAServerURL)
	if err != nil {
		return err
	}

	authz := authz.NewOPAClient(opaServerURL)

	h, err := handler.New(cfg, originURL, cfg.AccessToken, authz)
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

func generateBcrypt() error {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Printf("\nBcrypted: %s\n", string(b))
	return nil
}
