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

func DotFromPlanck(planck *big.Int) *big.Float {
	planckPerDot := big.NewFloat(substrate.PlankPerDot)
	planckFloat := new(big.Float).SetInt(planck)
	return new(big.Float).Quo(planckFloat, planckPerDot)
}

func PlanckFromDot(dot *big.Float) *big.Int {
	planckPerDot := big.NewFloat(substrate.PlankPerDot)
	v, _ := new(big.Float).Mul(dot, planckPerDot).Int(nil)
	return v
}

func dotsFromPlancks(plancks []*big.Int) []*big.Float {
	dots := make([]*big.Float, len(plancks))
	for i, p := range plancks {
		dots[i] = DotFromPlanck(p)
	}
	return dots
}
