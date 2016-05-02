package stream

import (
	"encoding/binary"
	"io"
)

// Media Box (mdia - mandatory)
//
// Contained in : Track Box (trak)
//
// Status: decoded
//
// Contains all information about the media data.
type MdiaBox struct {
	Mdhd   *MdhdBox
	Minf   *MinfBox
	boxes  []Box
	header [8]byte
}

func DecodeMdia(r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MdiaBox{
		boxes: make([]Box, 0, len(l)),
	}
	for _, b := range l {
		switch b.Type() {
		case "mdhd":
			m.Mdhd = b.(*MdhdBox)
		case "minf":
			m.Minf = b.(*MinfBox)
		default:
			m.boxes = append(m.boxes, b)
		}
	}
	return m, nil
}

func (b *MdiaBox) Type() string {
	return "mdia"
}

func (b *MdiaBox) Size() (sz int) {
	sz += b.Mdhd.Size()

	if b.Minf != nil {
		sz += b.Minf.Size()
	}

	for _, box := range b.boxes {
		sz += box.Size()
	}

	return sz + BoxHeaderSize
}

func (b *MdiaBox) Dump() {
	b.Mdhd.Dump()
	if b.Minf != nil {
		b.Minf.Dump()
	}
}

func (b *MdiaBox) Encode(w io.Writer) (err error) {
	binary.BigEndian.PutUint32(b.header[:4], uint32(b.Size()))
	copy(b.header[4:], b.Type())
	_, err = w.Write(b.header[:])
	if err != nil {
		return
	}
	err = b.Mdhd.Encode(w)
	if err != nil {
		return
	}

	for _, b := range b.boxes {
		if err = b.Encode(w); err != nil {
			return err
		}
	}

	return b.Minf.Encode(w)
}
