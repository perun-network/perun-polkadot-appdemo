package client

import (
	"context"
	"fmt"

	"github.com/perun-network/perun-polkadot-appdemo/app"
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
	_, err := func() (*pclient.LedgerChannelProposalMsg, error) {
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
		stake := lcp.FundingAgreement[assetIdx][ourIdx]
		if proposerStake.Cmp(stake) != 0 {
			return nil, fmt.Errorf("unequal stake")
		}

		// Propose to user.
		proposer := lcp.Peers[proposerIdx]
		stakeDot := DotFromPlanck(stake)
		msg := fmt.Sprintf("Incoming game proposal: Player %v, stake = %v DOT.\nAccept? (accept/reject)", proposer, stakeDot)
		h.io.Print(msg)
		h.proposals <- Proposal{lcp, r}

		return lcp, nil
	}()
	if err != nil {
		h.io.Print(fmt.Sprintf("Rejecting channel proposal: %v\n", err))
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
		return
	}
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
	p := func(msg string) {
		h.io.Print(fmt.Sprintf("Received event: %v (GameID: %x)", msg, e.ID()))
	}
	switch e := e.(type) {
	case *channel.RegisteredEvent:
		p("Dispute registered")
	case *channel.ProgressedEvent:
		p("State progressed")
		d, ok := e.State.Data.(*app.TicTacToeAppData)
		if !ok {
			h.io.Print(fmt.Sprintf("Error reading app state: wrong type: expected *TicTacToeAppData, got %T", e.State.Data))
			break
		}
		h.io.Print(fmt.Sprintf("New game state:\n%v", d))
	case *channel.ConcludedEvent:
		p("Concluded")
	default:
		msg := fmt.Sprintf("Unkown type %T", e)
		p(msg)
	}

}
