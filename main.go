package main

import (
	"log"
	"os"
	"github.com/ryotarai/github-api-authz-proxy/pkg/cli"
)

func main() {
	err := cli.New().Start(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}