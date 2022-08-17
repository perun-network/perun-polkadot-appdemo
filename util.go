package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net/simple"
)

func parseBigFloat(s string) (*big.Float, error) {
	v, ok := new(big.Float).SetString(s)
	if !ok {
		return nil, fmt.Errorf("error parsing big.Float: %v", s)
	}
	return v, nil
}

func parseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing int64: %v", err.Error()))
	}
	return i
}

func parseWireAddress(hex string) (wire.Address, error) {
	return simple.NewAddress(hex), nil
}

func hexToAddress(hex string) (*dotwallet.Address, error) {
	pk, err := sr25519.NewPKFromHex(hex)
	if err != nil {
		return nil, fmt.Errorf("creating public key from hex: %w", err)
	}
	return dotwallet.NewAddressFromPK(pk), nil
}
