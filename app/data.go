package app

import (
	"bytes"
	"fmt"

	dotchannel "github.com/perun-network/perun-polkadot-backend/channel"
	"perun.network/go-perun/channel"
)

// TicTacToeAppData is the app data struct.
// Grid:
// 0 1 2
// 3 4 5
// 6 7 8
type TicTacToeAppData struct {
	NextActor uint8
	Grid      [9]FieldValue
}

var _ channel.Data = &TicTacToeAppData{}

func (d *TicTacToeAppData) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%v|%v|%v\n", d.Grid[0], d.Grid[1], d.Grid[2])
	fmt.Fprintf(&b, "%v|%v|%v\n", d.Grid[3], d.Grid[4], d.Grid[5])
	fmt.Fprintf(&b, "%v|%v|%v\n", d.Grid[6], d.Grid[7], d.Grid[8])

	if final, winner := d.CheckFinal(); final {
		if winner == nil {
			fmt.Fprint(&b, "It's a draw.")
		} else {
			fmt.Fprintf(&b, "Winner: Player %v", *winner+1)
		}
	} else {
		fmt.Fprintf(&b, "Next actor: Player %v", d.NextActor+1)
	}
	return b.String()
}

// MarshalBinary encodes the data to bytes.
func (d *TicTacToeAppData) MarshalBinary() ([]byte, error) {
	return dotchannel.ScaleEncode(d)
}

// UnmarshalBinary decodes channel data from bytes.
func (d *TicTacToeAppData) UnmarshalBinary(data []byte) error {
	return dotchannel.ScaleDecode(d, data)
}

// Clone returns a deep copy of the app data.
func (d *TicTacToeAppData) Clone() channel.Data {
	_d := *d
	return &_d
}

func (d *TicTacToeAppData) Set(row, col int, actorIdx channel.Index) error {
	if d.NextActor != uint8safe(uint16(actorIdx)) {
		return fmt.Errorf("invalid actor")
	}
	v := makeFieldValueFromPlayerIdx(actorIdx)
	i := row*3 + col
	if i < 0 || i >= len(d.Grid) {
		return fmt.Errorf("index out of bounds: %d", i)
	}
	d.Grid[i] = v
	d.NextActor = calcNextActor(d.NextActor)
	return nil
}

func calcNextActor(actor uint8) uint8 {
	return (actor + 1) % numParts
}
