package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"perun.network/go-perun/wire"
)

func parseBigInt(s string) *big.Int {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("error parsing big.Int: %v", s))
	}
	return i
}

func parseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing int64: %v", err.Error()))
	}
	return i
}

func parseWireAddress(hex string) (wire.Address, error) {
	addr, err := hexToAddress(hex)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func hexToAddress(hex string) (*dotwallet.Address, error) {
	pk, err := sr25519.NewPKFromHex(hex)
	if err != nil {
		return nil, fmt.Errorf("creating public key from hex: %w", err)
	}
	return dotwallet.NewAddressFromPK(pk), nil
}
