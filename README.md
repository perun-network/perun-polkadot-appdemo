# Perun Polkadot App Demo

## Test

1. Start a local Polkadot node with the Perun Pallet and the Tic Tac Toe app installed at 127.0.0.1:9944.
You can do so by following the instructions at [perun-polkadot-node/wip-tictactoe](https://github.com/perun-network/perun-polkadot-node/tree/wip-tictactoe).

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


## TODO

- Investigate why event subscription is slow on sender side.
- Add CI with linter.
- Use interactive CI package (like gobata).
- Handle game and update in package main, not in package client. Communicate via channels in handler.
