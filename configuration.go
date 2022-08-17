package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net/simple"
)

type ConfigurationJSON struct {
	Host              string     `json:"host"`
	NodeURL           string     `json:"node_url"`
	NetworkID         uint8      `json:"network_id"`
	QueryDepth        uint32     `json:"query_depth"`
	DialTimeout       uint32     `json:"tx_timeout"`
	App               string     `json:"app_id"`
	ChallengeDuration uint64     `json:"challenge_duration"`
	SecretKey         string     `json:"secret_key"`
	Peers             []PeerJSON `json:"peers"`
}

type PeerJSON struct {
	Name        string `json:"name"`
	WireAddress string `json:"wire_address"`
	IpAddress   string `json:"ip_address"`
}

type Configuration struct {
	Host              string
	WireAccount       wire.Account
	NodeURL           string
	NetworkID         dot.NetworkID
	QueryDepth        types.BlockNumber
	DialTimeout       time.Duration
	App               AppID
	ChallengeDuration uint64
	SecretKey         string
	Peers             []Peer
}

type Peer struct {
	Name        string
	WireAddress wire.Address
	IpAddress   string
}

type AppID = wallet.Address

func loadConfig(fn string) (Configuration, error) {
	// Read file.
	text, err := ioutil.ReadFile(fn)
	if err != nil {
		return Configuration{}, err
	}

	// Unmarshal.
	var cfg ConfigurationJSON
	err = json.Unmarshal(text, &cfg)
	if err != nil {
		return Configuration{}, err
	}

	// Decode app.
	appAddr, err := hexToAddress(cfg.App)
	if err != nil {
		return Configuration{}, err
	}

	// Convert peers.
	peers := make([]Peer, len(cfg.Peers))
	for i, p := range cfg.Peers {
		peers[i] = makePeer(p)
	}

	return Configuration{
		Host:              cfg.Host,
		NodeURL:           cfg.NodeURL,
		NetworkID:         dot.NetworkID(cfg.NetworkID),
		QueryDepth:        types.BlockNumber(cfg.QueryDepth),
		DialTimeout:       time.Duration(cfg.DialTimeout) * time.Second,
		App:               appAddr,
		ChallengeDuration: cfg.ChallengeDuration,
		SecretKey:         cfg.SecretKey,
		Peers:             peers,
		WireAccount:       simple.NewAccount(simple.NewAddress(cfg.Host)),
	}, nil
}

func makePeer(p PeerJSON) Peer {
	addr, err := parseWireAddress(p.WireAddress)
	if err != nil {
		panic(err)
	}
	return Peer{
		Name:        p.Name,
		WireAddress: addr,
		IpAddress:   p.IpAddress,
	}
}
