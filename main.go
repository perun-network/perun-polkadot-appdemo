package main

import (
	"flag"
	"log"
	"os"

	"github.com/perun-network/perun-polkadot-appdemo/app"
	"github.com/perun-network/perun-polkadot-appdemo/cli"
	"github.com/perun-network/perun-polkadot-appdemo/client"
	"perun.network/go-perun/channel"
)

var cfgFlag = flag.String("cfg", "", "Configuration file")
var skFlag = flag.String("sk", "", "Secret key")

func main() {
	//TODO remove
	// rng := rand.New(rand.NewSource(0))
	// sk, _ := sr25519.NewSKFromRng(rng)
	// pk := sk.Public().Encode()
	// println(hex.EncodeToString(pk[:]))
	// os.Exit(0)

	flag.Parse()

	cfg, err := loadConfig(*cfgFlag)
	if err != nil {
		log.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	app := app.NewTicTacToeApp(cfg.App)
	channel.RegisterApp(app)

	init := func(io cli.IO) error {
		c, err := client.NewClient(
			*skFlag,
			cfg.NodeURL,
			cfg.NetworkID,
			cfg.QueryDepth,
			cfg.DialTimeout,
			app,
			io,
		)
		if err != nil {
			return err
		}
		io.SetContextValue(ContextKeyClient, c)
		io.SetContextValue(ContextKeyAddressBook, make(AddressBook))
		return nil
	}
	err = cli.Run(init, commands)
	if err != nil {
		log.Printf("Error running CLI: %v\n", err)
		os.Exit(1)
	}
}
