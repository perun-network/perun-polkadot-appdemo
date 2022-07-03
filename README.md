# Perun Polkadot App Demo

## Test

1. Start local Polkadot node with Perun Pallet at 127.0.0.1:9944.
2. Start client Bob.
```sh
go run . -cfg bob.config.json -sk 0x398f0c28f98885e046333d4a41c19cee4c37368a9832c6502f6cfd182e2aef89
```
```
3. Start client Alice in a new terminal.
```sh
go run . -cfg alice.config.json -sk 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a
```
4. In Alice's terminal, add Bob to Alice's peers.
```sh
addpeer bob 0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48 127.0.0.1:5751
```
5. Propose game to Bob with stake 10.
```sh
propose bob 10
```
6. Place mark at (2, 2).
```sh
set 2 2
```


## TODO

- Add CI with linter.
- Use Scale codec for encoding.
- Use interactive CI package (like gobata).
- Handle game and update in package main, not in package client. Communicate via channels in handler.
- Use host address as wire address.
- Handle player not responding.
- Add peers to config.
