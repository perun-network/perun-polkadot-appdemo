package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"perun.network/go-perun/wallet"
)

type ConfigurationJSON struct {
	NodeURL     string `json:"node_url"`
	NetworkID   uint8  `json:"network_id"`
	QueryDepth  uint32 `json:"query_depth"`
	DialTimeout uint32 `json:"tx_timeout"`
	App         string `json:"app_id"`
}

type Configuration struct {
	NodeURL     string
	NetworkID   dot.NetworkID
	QueryDepth  types.BlockNumber
	DialTimeout time.Duration
	App         AppID
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

	return Configuration{
		NodeURL:     cfg.NodeURL,
		NetworkID:   dot.NetworkID(cfg.NetworkID),
		QueryDepth:  types.BlockNumber(cfg.QueryDepth),
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
		App:         appAddr,
	}, nil
}
