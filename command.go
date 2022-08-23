package main

import (
	"fmt"
	"os"

	"github.com/perun-network/perun-polkadot-appdemo/cli"
	"github.com/perun-network/perun-polkadot-appdemo/client"
)

var commands = []cli.Command{
	{
		Name: "propose",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			if len(args) != 2 {
				io.Print("Invalid number of arguments.")
				return
			}
			peer := args[0]
			stake, err := parseBigFloat(args[1])
			if err != nil {
				io.Print(err.Error())
				return
			}
			stakePlanck := client.PlanckFromDot(stake)

			challengeDuration, err := Context(io).ChallengeDuration()
			if err != nil {
				io.Print(err.Error())
				return
			}

			// Get peer address.
			peerAddr, err := Context(io).PeerAddress(peer)
			if err != nil {
				io.Print(err.Error())
				return
			}

			// Get client.
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}

			// Propose game.
			ctx, cancel := c.NewTransactionContext()
			defer cancel()
			io.Print(fmt.Sprintf("Proposing game to %v (%v)...", peer, peerAddr))
			_, err = c.ProposeGame(ctx, peerAddr, stakePlanck, uint64(challengeDuration))
			if err != nil {
				io.Print("Error: " + err.Error())
				return
			}
		},
		Help: "Usage: propose [peer:string] [stake:float]\nPropose game to peer with stake in DOT.",
	},
	{
		Name: "accept",
		Func: func(io cli.IO, args []string) {
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			select {
			case p := <-c.Proposals():
				io.Print("Accepting proposal and depositing stake...")
				ctx, cancel := c.NewTransactionContext()
				defer cancel()
				err = c.AcceptProposal(ctx, p)
				if err != nil {
					io.Print(err.Error())
					return
				}
			default:
				io.Print("No incoming proposal.")
			}
		},
		Help: "Accept incoming proposal.",
	},
	{
		Name: "reject",
		Func: func(io cli.IO, args []string) {
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}

			select {
			case p := <-c.Proposals():
				io.Print("Rejecting proposal...")
				ctx, cancel := c.NewTransactionContext()
				defer cancel()
				err = c.RejectProposal(ctx, p, "rejected")
				if err != nil {
					io.Print(err.Error())
					return
				}
				io.Print("Done.")
			default:
				io.Print("No incoming proposal.")
			}
		},
		Help: "Reject incoming proposal.",
	},
	{
		Name: "set",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			expectedLen := 2
			if len(args) != expectedLen {
				io.Print(fmt.Sprintf("invalid number of arguments: expected %d, got %d", expectedLen, len(args)))
				return
			}
			row, column := parseInt64(args[0]), parseInt64(args[1])

			// Get game state.
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			g, err := c.Game()
			if err != nil {
				io.Print(err.Error())
				return
			}

			// Perform game action.
			io.Print(fmt.Sprintf("Proposing state update: Set mark at (%v, %v)", row, column))
			ctx, cancel := c.NewTransactionContext()
			defer cancel()
			err = g.Set(ctx, int(row)-1, int(column)-1)
			if err != nil {
				io.Print("Error performing game action: " + err.Error())
				return
			}
		},
		Help: "Usage: set [row:int] [column:int]\nSet mark.",
	},
	{
		Name: "force_set",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			expectedLen := 2
			if len(args) != expectedLen {
				io.Print(fmt.Sprintf("invalid number of arguments: expected %d, got %d", expectedLen, len(args)))
				return
			}
			row, column := parseInt64(args[0]), parseInt64(args[1])

			// Get game state.
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			g, err := c.Game()
			if err != nil {
				io.Print(err.Error())
				return
			}

			// Perform game action.
			io.Print(fmt.Sprintf("Forcing state update: Set mark at (%v, %v)", row, column))
			ctx, cancel := c.NewTransactionContext()
			defer cancel()
			err = g.ForceSet(ctx, int(row)-1, int(column)-1)
			if err != nil {
				io.Print("Error performing game action: " + err.Error())
				return
			}
			io.Print("Done.")
		},
		Help: "Usage: force_set [row:int] [column:int]\nEncforce action on-chain.",
	},
	{
		Name: "settle",
		Func: func(io cli.IO, args []string) {
			// Get game state.
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			g, err := c.Game()
			if err != nil {
				io.Print(err.Error())
				return
			}

			if !g.IsFinal() {
				io.Print("Error: Game not final.")
				return
			}

			// Close game.
			io.Print("Closing game...")
			ctx, cancel := c.NewTransactionContext()
			defer cancel()
			err = g.Settle(ctx)
			if err != nil {
				io.Print("Error closing game: " + err.Error())
				return
			}
			io.Print("Done.")
		},
		Help: "Settle the game state and withdraw funds. Errors if the game is not final.",
	},
	{
		Name: "force_settle",
		Func: func(io cli.IO, args []string) {
			// Get game state.
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			g, err := c.Game()
			if err != nil {
				io.Print(err.Error())
				return
			}

			// Close game.
			io.Print("Closing game...")
			ctx, cancel := c.NewTransactionContext()
			defer cancel()
			err = g.Settle(ctx)
			if err != nil {
				io.Print("Error closing game: " + err.Error())
				return
			}
			io.Print("Done.")
		},
		Help: "Settle the game state and withdraw funds. If the game is not final, game conclusion is being enforced.",
	},
	{
		Name: "state",
		Func: func(io cli.IO, args []string) {
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			g, err := c.Game()
			if err != nil {
				io.Print(err.Error())
				return
			}
			io.Print("Current game state:\n" + g.String())
		},
		Help: "Show my balance.",
	},
	{
		Name: "balance",
		Func: func(io cli.IO, args []string) {
			c, err := Context(io).Client()
			if err != nil {
				io.Print(err.Error())
				return
			}
			bal, err := c.Balance()
			if err != nil {
				io.Print(err.Error())
				return
			}
			dot := client.DotFromPlanck(bal.Int)
			io.Print("Balance: " + dot.String() + " DOT")
		},
		Help: "Show my balance.",
	},
	{
		Name: "addpeer",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			name := args[0]
			wireAddr, err := parseWireAddress(args[1])
			if err != nil {
				io.Print("Error parsing argument 2: " + err.Error())
				return
			}
			hostAddr := args[2]

			// Set peer address.
			err = Context(io).SetPeerAddress(name, wireAddr, hostAddr)
			if err != nil {
				io.Print("Error setting peer address: " + err.Error())
				return
			}

			io.Print(fmt.Sprintf("Added %v with wire address %v and network address %v to known peers.", name, wireAddr, hostAddr))
		},
		Help: "Usage: addpeer [name:string] [id:addr] [host:string]\nAdd client to list of known peers.",
	},
	{
		Name: "exit",
		Func: func(io cli.IO, args []string) {
			os.Exit(0)
		},
		Help: "Exit program.",
	},
}
