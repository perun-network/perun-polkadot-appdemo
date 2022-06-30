package client

import (
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/perun-network/perun-polkadot-appdemo/app"
	"github.com/perun-network/perun-polkadot-appdemo/cli"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/pkg/errors"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	wirenet "perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
	"perun.network/go-perun/wire/perunio/serializer"
)

type Client struct {
	perunClient *client.Client
	acc         wallet.Address
	app         *app.TicTacToeApp
	io          cli.IO
	game        *Game
	dialer      *simple.Dialer
}

func NewClient(
	sk string,
	nodeURL string,
	networkID dot.NetworkID,
	queryDepth types.BlockNumber,
	dialTimeout time.Duration,
	app *app.TicTacToeApp,
	io cli.IO,
) (*Client, error) {
	// Setup wallet.
	wallet, acc, err := setupWallet(sk)
	if err != nil {
		return nil, errors.WithMessage(err, "setting up wallet")
	}

	// Setup chain connection.
	dot, err := setupChain(acc, nodeURL, networkID, queryDepth)
	if err != nil {
		return nil, errors.WithMessage(err, "setting up chain connection")
	}

	// Setup network.
	dialer := simple.NewTCPDialer(dialTimeout)
	bus := wirenet.NewBus(acc, dialer, serializer.Serializer())

	// Setup watcher.
	watcher, err := local.NewWatcher(dot.Adjudicator)
	if err != nil {
		return nil, errors.WithMessage(err, "creating watcher")
	}

	// Setup Perun client.
	c, err := client.New(acc.Address(), bus, dot.Funder, dot.Adjudicator, wallet, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	gameClient := &Client{
		perunClient: c,
		acc:         acc.Address(),
		app:         app,
		io:          io,
		dialer:      dialer,
	}

	h := handler{gameClient}

	go c.Handle(h, h)

	return gameClient, nil
}
