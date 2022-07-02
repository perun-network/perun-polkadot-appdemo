package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/perun-network/perun-polkadot-appdemo/app"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

type Game struct {
	ch *client.Channel
}

func newGame(ch *client.Channel) *Game {
	return &Game{ch: ch}
}

// Set sends a game move to the channel peer.
func (g *Game) Set(row, col int) error {
	return g.ch.Update(context.TODO(), func(state *channel.State) {
		g.set(row, col, state)
	})
}

// ForceSet registers a game move on-chain.
func (g *Game) ForceSet(x, y int) error {
	return g.ch.ForceUpdate(context.TODO(), func(state *channel.State) {
		g.set(x, y, state)
	})
}

func (g *Game) set(row, col int, state *channel.State) {
	data, ok := state.Data.(*app.TicTacToeAppData)
	if !ok {
		panic(fmt.Sprintf("invalid data type: %T", data))
	}
	data.Set(row, col, g.ch.Idx())
	var winner *channel.Index
	state.IsFinal, winner = data.CheckFinal()
	if state.IsFinal {
		sum := state.Balances.Sum()
		if len(sum) != 1 {
			panic(fmt.Sprintf("expected 1 asset, got %v", len(sum)))
		}

		assetIdx := 0
		balances := state.Balances[assetIdx]
		for i, bal := range balances {
			if i == int(*winner) {
				bal.Set(sum[assetIdx])
			} else {
				bal.Set(big.NewInt(0))
			}
		}
	}
}

// Settle settles the app channel and withdraws the funds.
func (g *Game) Settle() error {
	err := g.ch.Settle(context.TODO(), false)
	if err != nil {
		return err
	}

	// Cleanup.
	g.ch.Close()
	return nil
}

func (g *Game) String() string {
	data, ok := g.ch.State().Data.(*app.TicTacToeAppData)
	if !ok {
		panic(fmt.Sprintf("invalid data type: %T", data))
	}
	return data.String()
}
