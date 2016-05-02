package stream

import (
	"encoding/binary"
	"io"
)

// Universal not decoded Box
type UniBox struct {
	name string
	buff []byte
	hbuf [BoxHeaderSize]byte
}

func DecodeUni(r io.Reader, name string) (Box, error) {
	data, err := readAllO(r)
	if err != nil {
		return nil, err
	}
	return &UniBox{
		name: name,
		buff: data,
	}, nil
}

func (b *UniBox) Type() string {
	return b.name
}

func (b *UniBox) Size() int {
	return BoxHeaderSize + len(b.buff)
}

func (b *UniBox) Encode(w io.Writer) (err error) {
	copy(b.hbuf[4:], b.Type())
	binary.BigEndian.PutUint32(b.hbuf[:4], uint32(b.Size()))

	if _, err = w.Write(b.hbuf[:]); err != nil {
		return
	}

	_, err = w.Write(b.buff)

	return
}
