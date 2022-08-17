package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/perun-network/perun-polkadot-appdemo/app"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

const assetIdx = 0

type Game struct {
	ch *client.Channel
}

func newGame(ch *client.Channel) *Game {
	return &Game{ch: ch}
}

// Set sends a game move to the channel peer.
func (g *Game) Set(ctx context.Context, row, col int) error {
	return g.applyAction(ctx, g.ch.Update, row, col)
}

// ForceSet registers a game move on-chain.
func (g *Game) ForceSet(ctx context.Context, row, col int) error {
	return g.applyAction(ctx, g.ch.ForceUpdate, row, col)
}

type updaterFn = func(ctx context.Context, updater func(*channel.State)) error

func (g *Game) applyAction(ctx context.Context, uf updaterFn, row, col int) error {
	// Dry run.
	err := g.set(row, col, g.ch.State().Clone())
	if err != nil {
		return fmt.Errorf("invalid move: %w", err)
	}

	return uf(ctx, func(state *channel.State) {
		err := g.set(row, col, state)
		if err != nil {
			panic(err)
		}
	})
}

func (g *Game) set(row, col int, state *channel.State) error {
	data, ok := state.Data.(*app.TicTacToeAppData)
	if !ok {
		panic(fmt.Sprintf("invalid data type: %T", data))
	}

	err := data.Set(row, col, g.ch.Idx())
	if err != nil {
		return err
	}

	var winner *channel.Index
	state.IsFinal, winner = data.CheckFinal()
	if state.IsFinal {
		sum := state.Balances.Sum()
		if len(sum) != 1 {
			panic(fmt.Sprintf("expected 1 asset, got %v", len(sum)))
		}

		balances := state.Balances[assetIdx]
		for i, bal := range balances {
			if i == int(*winner) {
				bal.Set(sum[assetIdx])
			} else {
				bal.Set(big.NewInt(0))
			}
		}
	}
	return nil
}

// Settle settles the app channel and withdraws the funds.
func (g *Game) Settle(ctx context.Context) error {
	err := g.ch.Settle(ctx, false)
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
	balances := dotsFromPlancks(g.ch.State().Balances[assetIdx])
	return fmt.Sprintf("Game ID: %x\nGame state:\n%v\nBalances: %v", g.ch.ID(), data.String(), balances)
}
