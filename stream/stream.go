package stream

import (
	"io"
	"time"
)

// A MPEG-4 media
//
// A MPEG-4 media contains three main boxes :
//
//   ftyp : the file type box
//   moov : the movie box (meta-data)
//   mdat : the media data (chunks and samples)
//
// Other boxes can also be present (pdin, moof, mfra, free, ...), but are not decoded.
type MP4 struct {
	Moov  *MoovBox
	Mdat  *MdatBox
	boxes []Box
}

// Decode decodes a media from a Reader
func Decode(r io.Reader) (*MP4, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	v := &MP4{
		boxes: make([]Box, 0, len(l)),
	}
	for _, b := range l {
		switch b.Type() {
		case "moov":
			v.Moov = b.(*MoovBox)
		case "mdat":
			v.Mdat = b.(*MdatBox)
		default:
			v.boxes = append(v.boxes, b)
		}
	}
	return v, nil
}

// Dump displays some information about a media
func (m *MP4) Dump() {
	m.Moov.Dump()
}

// Boxes lists the top-level boxes from a media
func (m *MP4) Boxes() []Box {
	return m.boxes
}

// Encode encodes a media to a Writer
func (m *MP4) Encode(w io.Writer) (err error) {
	err = m.Moov.Encode(w)
	if err != nil {
		return err
	}
	for _, b := range m.boxes {
		err = b.Encode(w)
		if err != nil {
			return err
		}
	}
	return m.Mdat.Encode(w)
}

func (m *MP4) Size() (sz int) {
	sz += m.Moov.Size()
	sz += m.Mdat.Size()

	for _, b := range m.Boxes() {
		sz += b.Size()
	}

	return
}

func (m *MP4) Duration() time.Duration {
	return time.Second * time.Duration(m.Moov.Mvhd.Duration) / time.Duration(m.Moov.Mvhd.Timescale)
}
