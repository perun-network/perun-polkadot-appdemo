package app

import (
	"fmt"

	"github.com/pkg/errors"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

var _ channel.App = &TicTacToeApp{} // Assert type implements interface.

// TicTacToeApp is a channel app.
type TicTacToeApp struct {
	Addr wallet.Address
}

func NewTicTacToeApp(addr wallet.Address) *TicTacToeApp {
	return &TicTacToeApp{
		Addr: addr,
	}
}

func (a *TicTacToeApp) NewData() channel.Data {
	return &TicTacToeAppData{}
}

// Def returns the app address.
func (a *TicTacToeApp) Def() wallet.Address {
	return a.Addr
}

func (a *TicTacToeApp) InitData(firstActor channel.Index) *TicTacToeAppData {
	return &TicTacToeAppData{
		NextActor: uint8(firstActor),
	}
}

// ValidInit checks that the initial state is valid.
func (a *TicTacToeApp) ValidInit(p *channel.Params, s *channel.State) error {
	if len(p.Parts) != numParts {
		return fmt.Errorf("invalid number of participants: expected %d, got %d", numParts, len(p.Parts))
	}

	appData, ok := s.Data.(*TicTacToeAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", s.Data)
	}

	zero := TicTacToeAppData{}
	if appData.Grid != zero.Grid {
		return fmt.Errorf("invalid starting grid: %v", appData.Grid)
	}

	if s.IsFinal {
		return fmt.Errorf("must not be final")
	}

	if appData.NextActor >= numParts {
		return fmt.Errorf("invalid next actor: got %d, expected < %d", appData.NextActor, numParts)
	}
	return nil
}

// ValidTransition is called whenever the channel state transitions.
func (a *TicTacToeApp) ValidTransition(params *channel.Params, from, to *channel.State, actor channel.Index) error {
	err := channel.AssertAssetsEqual(from.Assets, to.Assets)
	if err != nil {
		return fmt.Errorf("Invalid assets: %v", err)
	}

	fromData, ok := from.Data.(*TicTacToeAppData)
	if !ok {
		panic(fmt.Sprintf("from state: invalid data type: %T", from.Data))
	}

	toData, ok := to.Data.(*TicTacToeAppData)
	if !ok {
		panic(fmt.Sprintf("to state: invalid data type: %T", from.Data))
	}

	// Check actor.
	if fromData.NextActor != uint8safe(uint16(actor)) {
		return fmt.Errorf("invalid actor: expected %v, got %v", fromData.NextActor, actor)
	}

	// Check next actor.
	if len(params.Parts) != numParts {
		panic("invalid number of participants")
	}
	expectedToNextActor := calcNextActor(fromData.NextActor)
	if toData.NextActor != expectedToNextActor {
		return fmt.Errorf("invalid next actor: expected %v, got %v", expectedToNextActor, toData.NextActor)
	}

	// Check grid.
	changed := false
	for i, v := range toData.Grid {
		if v > maxFieldValue {
			return fmt.Errorf("invalid grid value at index %d: %d", i, v)
		}
		vFrom := fromData.Grid[i]
		if v != vFrom {
			if vFrom != notSet {
				return fmt.Errorf("cannot overwrite field %d", i)
			}
			if changed {
				return fmt.Errorf("cannot change two fields")
			}
			changed = true
		}
	}
	if !changed {
		return fmt.Errorf("cannot skip turn")
	}

	// Check balances.
	isFinal, winner := toData.CheckFinal()
	if to.IsFinal != isFinal {
		return fmt.Errorf("final flag: expected %v, got %v", isFinal, to.IsFinal)
	}
	expectedAllocation := from.Allocation.Clone()
	expectedAllocation.Balances = computeNextBalances(from.Allocation.Balances, actor, winner)
	if err := expectedAllocation.Equal(&to.Allocation); err != nil {
		return errors.WithMessagef(err, "wrong allocation: expected %v, got %v", expectedAllocation, to.Allocation)
	}
	return nil
}

func (*TicTacToeApp) Set(s *channel.State, x, y int, actorIdx channel.Index) error {
	d, ok := s.Data.(*TicTacToeAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	err := d.Set(x, y, actorIdx)
	if err != nil {
		return err
	}

	isFinal, winner := d.CheckFinal()
	s.IsFinal = isFinal
	s.Balances = computeNextBalances(s.Balances, actorIdx, winner)
	return nil
}
