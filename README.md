# Perun App Demo on Polkadot

This repository demonstrates how to play [Tic Tac Toe](https://en.wikipedia.org/wiki/Tic-tac-toe) for [Polkadot tokens](https://polkadot.network) in real time using [Perun Channels](https://github.com/perun-network/perun-polkadot-backend).

## Repository structure

* `app`: This package contains the off-chain definition of the Tic Tac Toe app.
* `cli`: This package provides an abstraction for the CLI interface.
* `client`: This package holds the Tic Tac Toe game client.


## Test instructions

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
