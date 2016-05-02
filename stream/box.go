package stream

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	BoxHeaderSize = 8
)

var (
	ErrTruncatedHeader = errors.New("truncated header")
)

var decoders map[string]BoxDecoder

func init() {
	decoders = map[string]BoxDecoder{
		"moov": DecodeMoov,
		"mvhd": DecodeMvhd,
		"trak": DecodeTrak,
		"tkhd": DecodeTkhd,
		"mdia": DecodeMdia,
		"minf": DecodeMinf,
		"mdhd": DecodeMdhd,
		"stbl": DecodeStbl,
		"stco": DecodeStco,
		"stsc": DecodeStsc,
		"stsz": DecodeStsz,
		"ctts": DecodeCtts,
		"stts": DecodeStts,
		"stss": DecodeStss,
		"mdat": DecodeMdat,
	}
}

// A box
type Box interface {
	Size() int
	Type() string
	Encode(w io.Writer) error
}

type BoxDecoder func(r io.Reader) (Box, error)

// DecodeContainer decodes a container box
func DecodeContainer(r io.Reader) (l []Box, err error) {
	var b Box
	var ht string
	var hs uint32

	buf := make([]byte, BoxHeaderSize)

	for {
		n, err := r.Read(buf)

		if err != nil {
			if err == io.EOF {
				return l, nil
			} else {
				return nil, err
			}
		}

		if n != BoxHeaderSize {
			return nil, ErrTruncatedHeader
		}

		ht = string(buf[4:8])
		hs = binary.BigEndian.Uint32(buf[0:4])

		if d := decoders[ht]; d != nil {
			b, err = d(io.LimitReader(r, int64(hs-BoxHeaderSize)))
		} else {
			b, err = DecodeUni(io.LimitReader(r, int64(hs-BoxHeaderSize)), ht)
		}

		if err != nil {
			return nil, err
		}

		l = append(l, b)

		if ht == "mdat" {
			b.(*MdatBox).ContentSize = hs - BoxHeaderSize
			return l, nil
		}
	}
}

// An 8.8 fixed point number
type Fixed16 uint16

func (f Fixed16) String() string {
	return fmt.Sprintf("%d.%d", uint16(f)>>8, uint16(f)&7)
}

func fixed16(bytes []byte) Fixed16 {
	return Fixed16(binary.BigEndian.Uint16(bytes))
}

func putFixed16(bytes []byte, i Fixed16) {
	binary.BigEndian.PutUint16(bytes, uint16(i))
}

// A 16.16 fixed point number
type Fixed32 uint32

func (f Fixed32) String() string {
	return fmt.Sprintf("%d.%d", uint32(f)>>16, uint32(f)&15)
}

func fixed32(bytes []byte) Fixed32 {
	return Fixed32(binary.BigEndian.Uint32(bytes))
}

func putFixed32(bytes []byte, i Fixed32) {
	binary.BigEndian.PutUint32(bytes, uint32(i))
}

// Utils
func makebuf(b Box) []byte {
	return make([]byte, b.Size()-BoxHeaderSize)
}

func readAllO(r io.Reader) ([]byte, error) {
	if lr, ok := r.(*io.LimitedReader); ok {
		buf := make([]byte, lr.N)
		_, err := io.ReadFull(lr, buf)
		return buf, err
	}
	return ioutil.ReadAll(r)
}
