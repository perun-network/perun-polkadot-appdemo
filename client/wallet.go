package client

import (
	"math/big"

	"github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"github.com/pkg/errors"
)

func setupWallet(hexSk string) (*dotwallet.Wallet, *dotwallet.Account, error) {
	wallet := dotwallet.NewWallet()
	sk, err := sr25519.NewSKFromHex(hexSk)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "creating hdwallet")
	}
	return wallet, wallet.ImportSK(sk), nil
}

func DotFromPlank(plank *big.Int) *big.Float {
	plankPerDot := big.NewFloat(substrate.PlankPerDot)
	plankFloat := new(big.Float).SetInt(plank)
	return new(big.Float).Quo(plankFloat, plankPerDot)
}

func PlankFromDot(dot *big.Float) *big.Int {
	plankPerDot := big.NewFloat(substrate.PlankPerDot)
	v, _ := new(big.Float).Mul(dot, plankPerDot).Int(nil)
	return v
}

func dotsFromPlanks(planks []*big.Int) []*big.Float {
	dots := make([]*big.Float, len(planks))
	for i, p := range planks {
		dots[i] = DotFromPlank(p)
	}
	return dots
}
