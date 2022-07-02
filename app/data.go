package app

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
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
		fmt.Fprintf(&b, "Winner: %v\n", *winner)
	} else {
		fmt.Fprintf(&b, "Next actor: %v\n", d.NextActor)
	}
	return b.String()
}

// MarshalBinary encodes the data to bytes.
func (d *TicTacToeAppData) MarshalBinary() ([]byte, error) {
	w := &bytes.Buffer{}
	err := writeUInt8(w, d.NextActor)
	if err != nil {
		return nil, errors.WithMessage(err, "writing actor")
	}

	err = writeUInt8Array(w, makeUInt8Array(d.Grid[:]))
	if err != nil {
		return nil, errors.WithMessage(err, "writing grid")
	}

	return w.Bytes(), nil
}

// UnmarshalBinary decodes channel data from bytes.
func (d *TicTacToeAppData) UnmarshalBinary(data []byte) error {
	r := bytes.NewBuffer(data)

	var err error
	d.NextActor, err = readUInt8(r)
	if err != nil {
		return errors.WithMessage(err, "reading actor")
	}

	grid, err := readUInt8Array(r, len(d.Grid))
	if err != nil {
		return errors.WithMessage(err, "reading grid")
	}
	copy(d.Grid[:], makeFieldValueArray(grid))
	return nil
}

// Clone returns a deep copy of the app data.
func (d *TicTacToeAppData) Clone() channel.Data {
	_d := *d
	return &_d
}

func (d *TicTacToeAppData) Set(row, col int, actorIdx channel.Index) {
	if d.NextActor != uint8safe(uint16(actorIdx)) {
		panic("invalid actor")
	}
	v := makeFieldValueFromPlayerIdx(actorIdx)
	d.Grid[col*3+row] = v
	d.NextActor = calcNextActor(d.NextActor)
}

func calcNextActor(actor uint8) uint8 {
	return (actor + 1) % numParts
}
