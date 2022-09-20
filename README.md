# `perun-polkadot-appdemo`: Trustlessly playing Tic-Tac-Toe for Polkadot tokens in Realtime using Perun App Channels

This repository demonstrates a CLI client that uses [Perun Channels](https://github.com/hyperledger-labs/go-perun) for trustlessly playing [Tic-Tac-Toe](https://en.wikipedia.org/wiki/Tic-tac-toe) for [Polkadot tokens](https://polkadot.network) in real time.

```
> set 1 1
Proposing state update: Set mark at (1, 1)
Game state:
x| | 
------
 | | 
------
 | | 
Next actor: Player 2
Balances: [20 0]
Waiting for other player.
```

## Introduction

The repository contains a CLI client that uses [perun-polkadot-backend](https://github.com/perun-network/perun-polkadot-backend) and [go-perun](https://github.com/hyperledger-labs/go-perun) to realize the functionality.

* `app`: The off-chain definition of the Tic Tac Toe app.
* `cli`: The CLI interface.
* `client`: The Tic Tac Toe game client.

The client can be started using [Go](https://go.dev).
It takes as input a config file.
```
go run . -cfg alice.config.json
```

Once the client is started, an overview of the available commands can be printed using the `help` command.

```
> help
```

## Testrun instructions

In the following, we describe how to test the client with a local Polkadot node.

1. Start a local [Polkadot node  with the Perun Pallet](https://github.com/perun-network/perun-polkadot-node) and the Tic Tac Toe app installed.
For example, you can do so using [Docker](https://www.docker.com).
```
docker run --rm -it -p9944:9944 ghcr.io/perun-network/polkadot-test-node:0.3.0
```

2. Start client Alice in one terminal.
```
go run . -cfg alice.config.json
```

3. Then start client Bob in a second terminal.
```
go run . -cfg bob.config.json
```

4. Propose game to Bob with stake 10.
```
propose bob 10
```
5. In Bob's terminal, accept the proposal.
```
accept
```
6. In Alice's terminal, set the first mark, e.g., at `(2, 2)`.
```
set 2 2
```
7. In Bob's terminal, set the next mark, e.g., at `(1, 1)`.
```
set 1 1
```
8. Play until game is in final state.
9. Conclude the game and withdraw the outcome.
```
settle
```

### Dispute

If the other client isn't responding, you can enforce an action on-chain.
```
force_set 3 3
```
At any time, the game can be forced to an end. The funds will be payed out according to the most recent state.
```
force_settle
```

## Funding

This project has been supported through the German Ministry for Education and Research (BMBF) and the Web3 Foundation.

## Copyright

Copyright 2022 PolyCrypt GmbH.
Use of the source code is governed by the [Apache License, Version 2.0](LICENSE).
