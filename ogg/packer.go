package ogg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
)

const preSkip = 312 // libopus encoder lookahead in 48kHz samples (per RFC 7845)

type Packer struct {
	channelCount uint8
	sampleRate   uint32
	serialNo     uint32
	packetNo     int64
	granulePos   int64
	buffer       bytes.Buffer
	oggEncoder   *Encoder
}

func New(channelCount uint8, sampleRate uint32, serialNo uint32) (*Packer, error) {
	p := Packer{
		channelCount: channelCount,
		sampleRate:   sampleRate,
		serialNo:     serialNo,
		packetNo:     1,
	}

	if err := p.init(); err != nil {
		return nil, fmt.Errorf("init ogg packer: %w", err)
	}

	runtime.SetFinalizer(&p, (*Packer).Close)

	return &p, nil
}

// AddChunk adds an opus packet to the ogg stream.
// samplesCount is the total number of PCM samples (across all channels) in the packet.
func (p *Packer) AddChunk(data []byte, eos bool, samplesCount int) error {
	numSamplesPerChannel := samplesCount / int(p.channelCount)

	// Granule positions are always in 48kHz samples per RFC 7845.
	p.granulePos += int64(numSamplesPerChannel) * 48000 / int64(p.sampleRate)

	if err := p.sendPacketToOggStream(data, false, eos); err != nil {
		return fmt.Errorf("send header data to ogg stream: %w", err)
	}

	return nil
}

// ReadPages returns all buffered ogg pages and sets the EOS flag on the last page.
func (p *Packer) ReadPages() ([]byte, error) {
	b := p.buffer.Bytes()
	if len(b) == 0 {
		return nil, errors.New("received empty ogg data buffer")
	}

	if err := setEOS(b); err != nil {
		return nil, fmt.Errorf("set EOS: %w", err)
	}

	p.buffer.Reset()
	return b, nil
}

func (p *Packer) Close() {
	p.oggEncoder = nil
	p.buffer.Reset()

	runtime.SetFinalizer(&p, nil)
}

func (p *Packer) init() error {
	p.oggEncoder = NewEncoder(p.serialNo, &p.buffer)

	if err := p.addHeader(); err != nil {
		return fmt.Errorf("add header to ogg stream: %w", err)
	}

	if err := p.addTags(); err != nil {
		return fmt.Errorf("add tags packet: %w", err)
	}

	return nil
}

func (p *Packer) addHeader() error {
	header := header(p.channelCount, p.sampleRate)
	if err := p.sendPacketToOggStream(header, true, false); err != nil {
		return fmt.Errorf("send header data to ogg stream: %w", err)
	}

	return nil
}

func (p *Packer) addTags() error {
	vendor := []byte("go-ogg-packer")
	tags := make([]byte, 8+4+len(vendor)+4) // magic + vendor_len + vendor + comment_count
	copy(tags, []byte("OpusTags"))
	binary.LittleEndian.PutUint32(tags[8:12], uint32(len(vendor)))
	copy(tags[12:], vendor)
	binary.LittleEndian.PutUint32(tags[12+len(vendor):], 0) // 0 comments

	if err := p.sendPacketToOggStream(tags, false, false); err != nil {
		return fmt.Errorf("send header data to ogg stream: %w", err)
	}

	return nil
}

// sendPacketToOggStream sends data to ogg stream in ogg packet format
// bos - begin of stream flag
// eos - end of stream flag
func (p *Packer) sendPacketToOggStream(data []byte, bos bool, eos bool) error {
	if bos {
		if err := p.oggEncoder.EncodeBOS(p.granulePos, [][]byte{data}); err != nil {
			return fmt.Errorf("write begin of stream packets to ogg stream: %w", err)
		}
		return nil
	}
	if eos {
		if err := p.oggEncoder.EncodeEOS(p.granulePos, [][]byte{data}); err != nil {
			return fmt.Errorf("write end of stream packets to ogg stream: %w", err)
		}
		return nil
	}

	if err := p.oggEncoder.Encode(p.granulePos, [][]byte{data}); err != nil {
		return fmt.Errorf("write packets to ogg stream: %w", err)
	}

	return nil
}

func header(channelCount uint8, sampleRate uint32) []byte {
	header := make([]byte, 19)
	copy(header, []byte("OpusHead"))

	header[8] = 1 // version number
	header[9] = channelCount

	binary.LittleEndian.PutUint16(header[10:12], preSkip)
	binary.LittleEndian.PutUint32(header[12:16], sampleRate)
	binary.LittleEndian.PutUint16(header[16:18], 0)

	header[18] = 0

	return header
}

// setEOS ensures EOS flag is set on the last page.
func setEOS(b []byte) error {
	lastPageStart := -1
	for i := len(b) - 4; i >= 0; i-- {
		if string(b[i:i+4]) == "OggS" {
			lastPageStart = i
			break
		}
	}
	if lastPageStart < 0 {
		return errors.New("no OggS page found in buffer")
	}

	if b[lastPageStart+5]&EOS != 0 {
		return nil
	}

	b[lastPageStart+5] |= EOS
	// Recalculate CRC after modifying the header.
	nSegs := int(b[lastPageStart+26])
	segTable := b[lastPageStart+27 : lastPageStart+27+nSegs]
	bodySize := 0
	for _, s := range segTable {
		bodySize += int(s)
	}
	pageBytes := b[lastPageStart : lastPageStart+27+nSegs+bodySize]
	pageBytes[22] = 0
	pageBytes[23] = 0
	pageBytes[24] = 0
	pageBytes[25] = 0
	crc := crc32(pageBytes)
	binary.LittleEndian.PutUint32(pageBytes[22:26], crc)

	return nil
}
