package packer

import (
	"fmt"
	"time"

	"github.com/paveldroo/go-ogg-packer/ogg"
	"github.com/paveldroo/go-ogg-packer/opus"
)

const DefaultPCMChunkSize = 2048

type Packer struct {
	result           []byte
	opusEncoder      *opus.Encoder
	oggPacker        *ogg.Packer
	pcmBuffer        []int16
	frameSizeSamples int
}

func New(numChannels int, sampleRate int, serialNo uint32) (*Packer, error) {
	cfg := opus.Config{
		SampleRate:  sampleRate,
		NumChannels: numChannels,
		FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
	}

	encoder, err := opus.NewEncoder(cfg)
	if err != nil {
		return nil, fmt.Errorf("create opus encoder: %s", err)
	}

	packer, err := ogg.New(uint8(numChannels), uint32(sampleRate), serialNo)
	if err != nil {
		return nil, fmt.Errorf("create ogg packer: %w", err)
	}

	return &Packer{
		opusEncoder:      encoder,
		oggPacker:        packer,
		frameSizeSamples: opus.FrameSizeSamples(cfg),
	}, nil
}

// SendPCMChunk encodes PCM data and packs it into OGG.
// Chunk must be PCM 16-bit little-endian audio samples.
// Chunk size is arbitrary — the packer buffers internally, so callers
// don't need to align chunks to frame boundaries.
func (s *Packer) SendPCMChunk(chunk []int16) error {
	s.pcmBuffer = append(s.pcmBuffer, chunk...)
	currentOpusPackets, pos, err := s.opusEncoder.Encode(s.pcmBuffer)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	s.pcmBuffer = s.pcmBuffer[pos:]
	for _, opusPacket := range currentOpusPackets {
		if err := s.oggPacker.AddChunk(opusPacket, false, s.frameSizeSamples); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}
	return nil
}

// GetResult flushes all PCM data and returns the complete OGG file.
// For streaming, use the ogg package directly with ReadPages/FlushPages.
func (s *Packer) GetResult() ([]byte, error) {
	defer s.oggPacker.Close()

	if err := s.flushPCMBuffer(); err != nil {
		return nil, fmt.Errorf("flush buffer: %w", err)
	}

	oggPages, err := s.oggPacker.FlushPages()
	if err != nil {
		return nil, fmt.Errorf("read pages: %w", err)
	}

	s.result = oggPages

	return s.result, nil
}

func (s *Packer) flushPCMBuffer() error {
	defer func() {
		s.pcmBuffer = s.pcmBuffer[:0]
	}()

	if len(s.pcmBuffer) == 0 {
		return nil
	}

	opusPackets, err := s.opusEncoder.EncodeWithPadding(s.pcmBuffer)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	for i, opusPacket := range opusPackets {
		eos := i == len(opusPackets)-1
		if err := s.oggPacker.AddChunk(opusPacket, eos, s.frameSizeSamples); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}

	return nil
}
