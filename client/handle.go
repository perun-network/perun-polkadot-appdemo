package client

import (
	"context"
	"fmt"

	dotchannel "github.com/perun-network/perun-polkadot-backend/channel"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	pclient "perun.network/go-perun/client"
)

type handler struct {
	*Client
}

// HandleProposal is the callback for incoming channel proposals.
func (h handler) HandleProposal(p pclient.ChannelProposal, r *pclient.ProposalResponder) {
	lcp, err := func() (*pclient.LedgerChannelProposalMsg, error) {
		// Check that we got a ledger channel proposal.
		lcp, ok := p.(*pclient.LedgerChannelProposalMsg)
		if !ok {
			return nil, fmt.Errorf("Invalid proposal type: %T\n", p)
		}

		// Check that the proposal includes the expected app.
		if !lcp.App.Def().Equal(h.app.Def()) {
			return nil, fmt.Errorf("Invalid appID: expected %v, got %v", h.app.Def(), lcp.App.Def())
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 {
			return nil, fmt.Errorf("Invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets.
		currency := dotchannel.Asset
		err := channel.AssertAssetsEqual(lcp.InitBals.Assets, []channel.Asset{currency})
		if err != nil {
			return nil, fmt.Errorf("Invalid assets: %v\n", err)
		}

		// Check that the funding agreement equals the initial balances.
		err = lcp.InitBals.Balances.AssertEqual(lcp.FundingAgreement)
		if err != nil {
			return nil, fmt.Errorf("funding agreement unequal to initial balances: %v\n", err)
		}

		// Check that the stakes are the same for all participants.
		const assetIdx, proposerIdx, ourIdx = 0, 0, 1
		proposerStake := lcp.FundingAgreement[assetIdx][proposerIdx]
		ourStake := lcp.FundingAgreement[assetIdx][ourIdx]
		if proposerStake.Cmp(ourStake) != 0 {
			return nil, fmt.Errorf("unequal stake")
		}

		// Propose to user.
		proposer := lcp.Peers[proposerIdx]
		msg := fmt.Sprintf("Incoming game proposal: Player %v, stake = %v. Accept? (y/n)", proposer, ourStake)
		answer, err := h.io.Prompt(msg)
		if err != nil {
			return nil, fmt.Errorf("prompting user input")
		} else if answer != "y" {
			return nil, fmt.Errorf("proposal rejected")
		}

		return lcp, nil
	}()
	if err != nil {
		h.io.Print(fmt.Sprintf("Rejecting channel proposal: %v\n", err))
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
		return
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		h.acc,                    // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := r.Accept(context.TODO(), accept)
	if err != nil {
		h.io.Print(fmt.Sprintf("Error accepting channel proposal: %v\n", err))
		return
	}

	// Start the on-chain event watcher.
	go func() {
		err := ch.Watch(h)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()

	h.game = newGame(ch)
}

// HandleUpdate is the callback for incoming channel updates.
func (h handler) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// Perun automatically checks that the transition is valid.
	// We always accept.
	err := r.Accept(context.TODO())
	if err != nil {
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (h handler) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	h.io.Print(fmt.Sprintf("Received adjudicator event: %T", e))
}
