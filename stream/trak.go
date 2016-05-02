package stream

import (
	"encoding/binary"
	"io"
)

// Track Box (tkhd - mandatory)
//
// Contained in : Movie Box (moov)
//
// A media file can contain one or more tracks.
type TrakBox struct {
	Tkhd   *TkhdBox
	Mdia   *MdiaBox
	boxes  []Box
	header [8]byte
}

func DecodeTrak(r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	t := &TrakBox{
		boxes: make([]Box, 0, len(l)),
	}
	for _, b := range l {
		switch b.Type() {
		case "tkhd":
			t.Tkhd = b.(*TkhdBox)
		case "mdia":
			t.Mdia = b.(*MdiaBox)
		default:
			t.boxes = append(t.boxes, b)
		}
	}
	return t, nil
}

func (b *TrakBox) Type() string {
	return "trak"
}

func (b *TrakBox) Size() (sz int) {
	sz += b.Tkhd.Size()
	sz += b.Mdia.Size()

	for _, box := range b.boxes {
		sz += box.Size()
	}

	return sz + BoxHeaderSize
}

func (b *TrakBox) Dump() {
	b.Tkhd.Dump()
	b.Mdia.Dump()
}

func (b *TrakBox) Encode(w io.Writer) (err error) {
	binary.BigEndian.PutUint32(b.header[:4], uint32(b.Size()))
	copy(b.header[4:], b.Type())
	_, err = w.Write(b.header[:])
	if err != nil {
		return
	}

	if err = b.Tkhd.Encode(w); err != nil {
		return
	}

	for _, b := range b.boxes {
		if err = b.Encode(w); err != nil {
			return
		}
	}

	return b.Mdia.Encode(w)
}
