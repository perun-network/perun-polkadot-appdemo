package app

import (
	"fmt"
	"math/big"

	"perun.network/go-perun/channel"
)

const numParts = 2

type FieldValue uint8

const (
	notSet FieldValue = iota
	player1
	player2
	maxFieldValue = player2
)

func (v FieldValue) String() string {
	switch v {
	case notSet:
		return " "
	case player1:
		return "x"
	case player2:
		return "o"
	default:
		panic(fmt.Sprintf("unsupported value: %d", v))
	}
}

func makeFieldValueFromPlayerIdx(idx channel.Index) FieldValue {
	switch idx {
	case 0:
		return player1
	case 1:
		return player2
	default:
		panic("invalid")
	}
}

func (v FieldValue) PlayerIndex() channel.Index {
	switch v {
	case player1:
		return 0
	case player2:
		return 1
	default:
		panic("invalid")
	}
}

func (d TicTacToeAppData) CheckFinal() (isFinal bool, winner *channel.Index) {
	// 0 1 2
	// 3 4 5
	// 6 7 8

	// Check winner.
	v := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // columns
		{0, 4, 8}, {2, 4, 6}, // diagonals
	}

	for _, _v := range v {
		ok, idx := d.samePlayer(_v...)
		if ok {
			return true, &idx
		}
	}

	// Check all set.
	for _, v := range d.Grid {
		if v == notSet {
			return false, nil
		}
	}
	return true, nil
}

func (d TicTacToeAppData) samePlayer(gridIndices ...int) (ok bool, player channel.Index) {
	if len(gridIndices) < 2 {
		panic("expecting at least two inputs")
	}

	first := d.Grid[gridIndices[0]]
	if first == notSet {
		return false, 0
	}
	for _, i := range gridIndices {
		if d.Grid[i] != first {
			return false, 0
		}
	}
	return true, first.PlayerIndex()
}

func uint8safe(a uint16) uint8 {
	b := uint8(a)
	if uint16(b) != a {
		panic("unsafe")
	}
	return b
}

func computeNextBalances(bals channel.Balances, actor channel.Index, winner *channel.Index) channel.Balances {
	total := bals.Sum()
	nextBals := bals.Clone()
	for a, assetBals := range nextBals {
		for p, b := range assetBals {
			p := channel.Index(p)
			if winner != nil && *winner == p || actor == p {
				b.Set(total[a])
			} else {
				b.Set(big.NewInt(0))
			}
		}
	}
	return nextBals
}
