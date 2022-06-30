package client

import (
	"github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
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
