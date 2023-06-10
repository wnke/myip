package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/wnke/myip"
)

func main() {
	app := &cli.App{
		Name:  "myip",
		Usage: "Get the your public IP address!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required:    false,
				Name:        "provider",
				Usage:       "IP provider URL like https://ifconfig.me, if empty a random one will be used",
				DefaultText: "",
				Value:       "",
				EnvVars:     []string{"MYIP_PROVIDER"},
			},
		},

		Action: func(ctx *cli.Context) error {
			providerArg := ctx.String("provider")
			var providers []string
			if providerArg != "" {
				providers = []string{providerArg}
			} else {
				providers = []string{myip.DEFAULT_PROVIDERS[rand.Intn(len(myip.DEFAULT_PROVIDERS))]}
			}

			ipa, err := myip.NewIPDiscoverWithProviders(providers)
			if err != nil {
				return err
			}

			ip, err := ipa.Discover()
			if err != nil {
				return err
			}
			fmt.Print(ip.String())

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
