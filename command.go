package main

import (
	"github.com/perun-network/perun-polkadot-appdemo/cli"
	"github.com/perun-network/perun-polkadot-appdemo/client"
)

var commands = []cli.Command{
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
		},
		Help: "addpeer [name:string] [id:addr] [host:string]",
	},
	{
		Name: "propose",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			if len(args) != 3 {
				io.Print("Invalid number of arguments.")
				return
			}
			peer, stake, challengeDuration := args[0], parseBigInt(args[1]), parseInt64(args[2])

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
			g, err := c.ProposeGame(peerAddr, stake, uint64(challengeDuration))
			if err != nil {
				io.Print(err.Error())
				return
			}
			io.SetContextValue(ContextKeyGame, g)
			io.Print(g.String())
		},
		Help: "propose [peer:addr] [stake:int] [challengeDuration:int]",
	},
	{
		Name: "set",
		Func: func(io cli.IO, args []string) {
			// Parse arguments.
			x, y := parseInt64(args[0]), parseInt64(args[1])

			// Get game state.
			gUntyped, ok := io.ContextValue(ContextKeyGame)
			if !ok {
				io.Print("Game not initialized.")
				return
			}
			g := gUntyped.(*client.Game)

			// Perform game action.
			err := g.Set(int(x), int(y))
			if err != nil {
				io.Print("Error performing game action: " + err.Error())
				return
			}
		},
		Help: "set [x:int] [y:int]",
	},
}
