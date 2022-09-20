package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/perun-network/perun-polkadot-appdemo/app"
	"github.com/perun-network/perun-polkadot-appdemo/cli"
	"github.com/perun-network/perun-polkadot-appdemo/client"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"perun.network/go-perun/channel"
)

var cfgFlag = flag.String("cfg", "", "Configuration file")

func main() {
	flag.Parse()

	cfg, err := loadConfig(*cfgFlag)
	if err != nil {
		log.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	app := app.NewTicTacToeApp(cfg.App)
	channel.RegisterApp(app)

	init := func(io *cli.IO) error {
		api, err := substrate.NewAPI(cfg.NodeURL, cfg.NetworkID)
		if err != nil {
			return err
		}

		c, err := client.NewClient(
			cfg.SecretKey,
			api,
			cfg.QueryDepth,
			cfg.Host,
			cfg.WireAccount,
			cfg.TxTimeout,
			app,
			io,
		)
		if err != nil {
			return err
		}
		io.SetContextValue(ContextKeyClient, c)
		io.SetContextValue(ContextKeyChallengeDuration, cfg.ChallengeDuration)

		io.Print(fmt.Sprintf("Account: %v", c.Account()))
		bal, err := c.Balance()
		if err != nil {
			return err
		}
		balDot := client.DotFromPlanck(bal.Int)
		io.Print(fmt.Sprintf("Balance: %v DOT", balDot.String()))

		// Add Peers.
		book := make(AddressBook)
		for _, p := range cfg.Peers {
			c.RegisterPeer(p.WireAddress, p.IpAddress)
			book[p.Name] = p.WireAddress
		}
		io.SetContextValue(ContextKeyAddressBook, book)
		return nil
	}
	err = cli.Run(init, commands)
	if err != nil {
		log.Printf("Error running CLI: %v\n", err)
		os.Exit(1)
	}
}
