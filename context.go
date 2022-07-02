package main

import (
	"fmt"

	"github.com/perun-network/perun-polkadot-appdemo/cli"
	"github.com/perun-network/perun-polkadot-appdemo/client"
	"perun.network/go-perun/wire"
)

const (
	ContextKeyClient            = "Client"
	ContextKeyAddressBook       = "AddressBook"
	ContextKeyChallengeDuration = "ChallengeDuration"
)

type Context cli.IO

func (c Context) Client() (*client.Client, error) {
	cIface, ok := cli.IO(c).ContextValue(ContextKeyClient)
	if !ok {
		return nil, fmt.Errorf("could not load client")
	}

	client, ok := cIface.(*client.Client)
	if !ok {
		return nil, fmt.Errorf("wrong type")
	}

	return client, nil
}

type AddressBook map[string]wire.Address

func (c Context) AddressBook() (AddressBook, error) {
	addrBookIface, ok := cli.IO(c).ContextValue(ContextKeyAddressBook)
	if !ok {
		return nil, fmt.Errorf("could not load address book")
	}

	addrBook, ok := addrBookIface.(AddressBook)
	if !ok {
		return nil, fmt.Errorf("wrong type")
	}
	return addrBook, nil
}

func (c Context) PeerAddress(peer string) (wire.Address, error) {
	// Load address book.
	addrBook, err := c.AddressBook()
	if err != nil {
		return nil, err
	}

	// Resolve name.
	peerAddr, ok := addrBook[peer]
	if !ok {
		return nil, fmt.Errorf("could not find address for peer: %v", peer)
	}

	return peerAddr, nil
}

func (c Context) SetPeerAddress(name string, wireAddr wire.Address, hostAddr string) error {
	// Load address book.
	addrBook, err := c.AddressBook()
	if err != nil {
		return err
	}

	// Register peer on wire dialer.
	client, err := c.Client()
	if err != nil {
		return err
	}
	client.RegisterPeer(wireAddr, hostAddr)

	// Set peer address.
	addrBook[name] = wireAddr

	return nil
}

func (c Context) ChallengeDuration() (uint64, error) {
	durationIface, ok := cli.IO(c).ContextValue(ContextKeyChallengeDuration)
	if !ok {
		return 0, fmt.Errorf("could not load challenge duration")
	}

	duration, ok := durationIface.(uint64)
	if !ok {
		return 0, fmt.Errorf("wrong type")
	}
	return duration, nil
}
