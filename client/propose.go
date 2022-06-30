package client

import (
	"context"
	"fmt"

	dotchannel "github.com/perun-network/perun-polkadot-backend/channel"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wire"
)

func (c *Client) ProposeGame(peer wire.Address, stake channel.Bal, challengeDuration uint64) (*Game, error) {
	participants := []wire.Address{c.acc, peer}

	// We create an initial allocation which defines the starting balances.
	currency := dotchannel.Asset
	initAlloc := channel.NewAllocation(2, currency)
	initAlloc.SetAssetBalances(currency, []channel.Bal{
		stake, // Our initial balance.
		stake, // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	firstActorIdx := channel.Index(0)
	app := c.app
	appParam := client.WithApp(app, app.InitData(firstActorIdx))

	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.acc,
		initAlloc,
		participants,
		appParam,
	)
	if err != nil {
		return nil, err
	}

	// Send the app channel proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		return nil, err
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	go func() {
		h := handler{c}
		err := ch.Watch(h)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()

	return newGame(ch), nil
}
