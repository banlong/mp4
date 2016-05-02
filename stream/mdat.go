package stream

import (
	"encoding/binary"
	"io"
)

// Media Data Box (mdat - optional)
//
// Status: not decoded
//
// The mdat box contains media chunks/samples.
//
// It is not read, only the io.Reader is stored, and will be used to Encode (io.Copy) the box to a io.Writer.
type MdatBox struct {
	ContentSize uint32
	header      [8]byte
	r           io.Reader
}

func DecodeMdat(r io.Reader) (Box, error) {
	// r is a LimitedReader
	if lr, limited := r.(*io.LimitedReader); limited {
		r = lr.R
	}
	return &MdatBox{r: r}, nil
}

func (b *MdatBox) Type() string {
	return "mdat"
}

func (b *MdatBox) Size() int {
	return BoxHeaderSize + int(b.ContentSize)
}

func (b *MdatBox) Reader() io.Reader {
	return b.r
}

func (b *MdatBox) Encode(w io.Writer) error {
	binary.BigEndian.PutUint32(b.header[:4], uint32(b.Size()))
	copy(b.header[4:], b.Type())
	_, err := w.Write(b.header[:])
	if err != nil {
		return err
	}
	_, err = io.Copy(w, b.r)
	return err
}
