package main

import (
	"fmt"
	"os"

	"github.com/perun-network/perun-polkadot-appdemo/cli"
)

var commands = []cli.Command{
	// {
	// 	Name: "addpeer",
	// 	Func: func(io cli.IO, args []string) {
	// 		// Parse arguments.
	// 		name := args[0]
	// 		wireAddr, err := parseWireAddress(args[1])
	// 		if err != nil {
	// 			io.Print("Error parsing argument 2: " + err.Error())
	// 			return
	// 		}
	// 		hostAddr := args[2]

	// 		// Set peer address.
	// 		err = Context(io).SetPeerAddress(name, wireAddr, hostAddr)
	// 		if err != nil {
	// 			io.Print("Error setting peer address: " + err.Error())
	// 			return
	// 		}

	// 		io.Print(fmt.Sprintf("Added %v with wire address %v and network address %v to known peers.", name, wireAddr, hostAddr))
	// 	},
	// 	Help: "Usage: addpeer [name:string] [id:addr] [host:string]\nAdd client to list of known peers.",
	// },
	{
		Name: "propose",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			if len(args) != 2 {
				io.Print("Invalid number of arguments.")
				return
			}
			peer, stake := args[0], parseBigInt(args[1])

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
			io.Print(fmt.Sprintf("Proposing game to %v (%v)...", peer, peerAddr))
			_, err = c.ProposeGame(peerAddr, stake, uint64(challengeDuration))
			if err != nil {
				io.Print("Error: " + err.Error())
				return
			}
		},
		Help: "Usage: propose [peer:string] [stake:int]\nPropose game to peer.",
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
				err = c.AcceptProposal(p)
				if err != nil {
					io.Print(err.Error())
					return
				}
				io.Print("Done.")
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
				err = c.RejectProposal(p, "rejected")
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
		Name: "mark",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			expectedLen := 2
			if len(args) != expectedLen {
				io.Print(fmt.Sprintf("invalid number of arguments: expected %d, got %d", expectedLen, len(args)))
				return
			}
			row, column := parseInt64(args[0])-1, parseInt64(args[1])-1

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
			io.Print(fmt.Sprintf("Proposing state update: place mark at (%v, %v)", row, column))
			err = g.Set(int(row), int(column))
			if err != nil {
				io.Print("Error performing game action: " + err.Error())
				return
			}
			io.Print("Update accepted.")
		},
		Help: "Usage: mark [row:int] [column:int]\nPlace mark.",
	},
	{
		Name: "force_mark",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			expectedLen := 2
			if len(args) != expectedLen {
				io.Print(fmt.Sprintf("invalid number of arguments: expected %d, got %d", expectedLen, len(args)))
				return
			}
			row, column := parseInt64(args[0])-1, parseInt64(args[1])-1

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
			io.Print(fmt.Sprintf("Forcing state update: place mark at (%v, %v)", row, column))
			err = g.ForceSet(int(row), int(column))
			if err != nil {
				io.Print("Error performing game action: " + err.Error())
				return
			}
			io.Print("Done.")
		},
		Help: "Usage: force_mark [row:int] [column:int]\nEncforce action on-chain.",
	},
	{
		Name: "force_quit",
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
			err = g.Settle()
			if err != nil {
				io.Print("Error closing game: " + err.Error())
				return
			}
			io.Print("Done.")
		},
		Help: "Force the game to come to an end (e.g., if the other participant does not respond).",
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
			io.Print("Balance: " + bal.String())
		},
		Help: "Show my balance.",
	},
	{
		Name: "exit",
		Func: func(io cli.IO, args []string) {
			os.Exit(0)
		},
		Help: "Exit program.",
	},
}
