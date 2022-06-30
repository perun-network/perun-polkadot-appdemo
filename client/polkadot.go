package client

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"

	"github.com/perun-network/perun-polkadot-backend/channel/pallet"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	pchannel "perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"
)

type (
	chain struct {
		Api         *dot.API
		Funder      pchannel.Funder
		Adjudicator pchannel.Adjudicator
	}
)

func setupChain(
	acc pwallet.Account,
	api *dot.API,
	queryDepth types.BlockNumber,
) (*chain, error) {
	perun := pallet.NewPallet(pallet.NewPerunPallet(api), api.Metadata())
	funder := pallet.NewFunder(perun, acc, 3)
	adj := pallet.NewAdjudicator(acc, perun, api, queryDepth)
	return &chain{api, funder, adj}, nil
}
