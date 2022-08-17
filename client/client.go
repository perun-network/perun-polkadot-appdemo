package client

import (
	"context"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/perun-network/perun-polkadot-appdemo/app"
	"github.com/perun-network/perun-polkadot-appdemo/cli"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	pclient "perun.network/go-perun/client"
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
	api         *dot.API
	proposals   chan Proposal
}

func NewClient(
	sk string,
	api *dot.API,
	queryDepth types.BlockNumber,
	host string,
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
	dot, err := setupChain(acc, api, queryDepth)
	if err != nil {
		return nil, errors.WithMessage(err, "setting up chain connection")
	}

	// Setup network.
	dialer := simple.NewTCPDialer(dialTimeout)
	bus := wirenet.NewBus(acc, dialer, serializer.Serializer())
	listener, err := simple.NewTCPListener(host)
	if err != nil {
		return nil, errors.WithMessage(err, "could not start tcp listener")
	}

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
		api:         api,
		proposals:   make(chan Proposal),
	}

	h := handler{gameClient}

	go c.Handle(h, h)
	go bus.Listen(listener)

	return gameClient, nil
}

func (c *Client) Game() (*Game, error) {
	if c.game == nil {
		return nil, fmt.Errorf("game not set")
	}
	return c.game, nil
}

func (c *Client) Account() dotwallet.Address {
	return dotwallet.AsAddr(c.acc)
}

func (c *Client) Balance() (types.U128, error) {
	addr := dotwallet.AsAddr(c.acc)
	info, err := c.api.AccountInfo(addr.AccountID())
	if err != nil {
		return types.U128{}, err
	}
	return info.Free, nil
}

func (c *Client) initGame(ch *client.Channel) {
	// Start the on-chain event watcher.
	go func() {
		h := handler{c}
		err := ch.Watch(h)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()

	// Handle updates.
	ch.OnUpdate(func(from, to *channel.State) {
		data := to.Data.(*app.TicTacToeAppData)
		balances := dotsFromPlancks(to.Balances[assetIdx])
		c.io.Print(fmt.Sprintf("***\nUpdated game state:\n%v\nBalances: %v", data.String(), balances))

		// If final, settle.
		if final, winner := data.CheckFinal(); final {
			if c.game.ch.Idx() == *winner {
				c.io.Print("You won.")
				c.io.Print("Initiating payout...")
				go func() {
					err := ch.Settle(context.TODO(), false)
					if err != nil {
						c.io.Print(err.Error())
						return
					}
					c.io.Print("Payout done.")
					ch.Close()
				}()
			} else {
				c.io.Print("Game over. You lost.")
			}
		} else {
			if c.game.ch.Idx() == channel.Index(data.NextActor) {
				c.io.Print("Your turn.")
			} else {
				c.io.Print("Waiting for other player.")
			}
		}
	})

	c.game = newGame(ch)
	c.io.Print("New game started.\n" + c.game.String())
	data := ch.State().Data.(*app.TicTacToeAppData)
	if c.game.ch.Idx() == channel.Index(data.NextActor) {
		c.io.Print("Your turn.")
	} else {
		c.io.Print("Waiting for other player.")
	}
}

type Proposal struct {
	p *pclient.LedgerChannelProposalMsg
	r *pclient.ProposalResponder
}

func (c *Client) Proposals() chan Proposal {
	return c.proposals
}

func (c *Client) AcceptProposal(p Proposal) error {
	// Create a channel accept message and send it.
	accept := p.p.Accept(
		c.acc,                    // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := p.r.Accept(context.TODO(), accept)
	if err != nil {
		return err
	}
	c.initGame(ch)
	return nil
}

func (c *Client) RejectProposal(p Proposal, reason string) error {
	return p.r.Reject(context.TODO(), reason)
}
