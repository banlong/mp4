package stream

import (
	"encoding/binary"
	"io"
)

// Sample Size Box (stsz - mandatory)
//
// Contained in : Sample Table box (stbl)
//
// Status : decoded
//
// For each track, either stsz of the more compact stz2 must be present. stz2 variant is not supported.
//
// This table lists the size of each sample. If all samples have the same size, it can be defined in the
// SampleUniformSize attribute.
type StszBox struct {
	body   []byte
	header [8]byte

	SampleStart       uint32
	SampleNumber      uint32
	SampleUniformSize uint32
}

func DecodeStsz(r io.Reader) (Box, error) {
	data, err := readAllO(r)

	if err != nil {
		return nil, err
	}

	b := &StszBox{
		body:              data,
		SampleUniformSize: binary.BigEndian.Uint32(data[4:8]),
	}

	return b, nil
}

func (b *StszBox) Type() string {
	return "stsz"
}

func (b *StszBox) Size() int {
	return BoxHeaderSize + 12 + int(b.SampleNumber)*4
}

func (b *StszBox) Encode(w io.Writer) (err error) {
	defer func() {
		b.body = nil
	}()

	binary.BigEndian.PutUint32(b.header[:4], uint32(b.Size()))
	copy(b.header[4:], b.Type())

	if _, err = w.Write(b.header[:]); err != nil {
		return
	}

	binary.BigEndian.PutUint32(b.body[8:12], uint32(b.SampleNumber))

	if _, err = w.Write(b.body[:12]); err != nil {
		return
	}

	if b.SampleUniformSize == 0 {
		if _, err = w.Write(b.body[12+4*b.SampleStart : 16+4*(b.SampleStart+b.SampleNumber-1)]); err != nil {
			return
		}
	}

	return err
}

// GetSampleSize returns the size (in bytes) of a sample
func (b *StszBox) GetSampleSize(i int) uint32 {
	if b.SampleUniformSize > 0 {
		return b.SampleUniformSize
	}

	return binary.BigEndian.Uint32(b.body[(12 + 4*i):(16 + 4*i)])
}
