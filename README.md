# Perun Polkadot App Demo

## Test

1. Start local Polkadot node with Perun Pallet at 127.0.0.1:9944.

```sh
 docker run --rm -p 9944:9944 ghcr.io/perun-network/polkadot-test-node:0.2.0
```

2. Start client Alice in one terminal.
```sh
go run . -cfg alice.config.json
```

3. Then start client Bob in a second terminal.
```sh
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
6. In Alice's terminal, place the first mark, e.g., at `(2, 2)`.
```
mark 2 2
```
7. In Bob's terminal, place second next mark, e.g., at `(1, 1)`.
```
mark 1 1
```
8. Play until game is in final state.
9. Conclude the game and withdraw funds.
```
conclude
```

### Dispute

If the other client isn't responding, you can enforce an action on-chain.
```
force_mark 3 3
```
At any time, the game can be forced to an end. The funds will be payed out according to the most recent state.
```
force_conclude
```


## TODO

- Investigate why event subscription is slow on sender side.
- Add CI with linter.
- Use interactive CI package (like gobata).
- Handle game and update in package main, not in package client. Communicate via channels in handler.
