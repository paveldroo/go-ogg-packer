package writer

import (
	"fmt"

	packer "github.com/paveldroo/go-ogg-packer"
	"github.com/paveldroo/go-ogg-packer/tests/buffer_writer/opus_tools"
)

type AudioBufferWriter struct {
	result        []byte
	opusConverter *opus_tools.Converter
	oggPacker     *packer.Packer
	lastS16Buffer []int16
}

func NewAudioBuffer(
	opusConverter *opus_tools.Converter,
	oggPacker *packer.Packer,
) *AudioBufferWriter {
	return &AudioBufferWriter{
		opusConverter: opusConverter,
		oggPacker:     oggPacker,
	}
}

func (s *AudioBufferWriter) SendS16Chunk(chunk []int16) error {
	s.lastS16Buffer = append(s.lastS16Buffer, chunk...)
	currentOpusPackets, pos, err := s.opusConverter.Encode(s.lastS16Buffer)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	s.lastS16Buffer = s.lastS16Buffer[pos:]
	for _, opusPacket := range currentOpusPackets {
		if err := s.oggPacker.AddChunk(opusPacket, false, pos); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}
	return nil
}

func (s *AudioBufferWriter) GetResult() ([]byte, error) {
	defer s.oggPacker.Close()

	if err := s.flushLastS16Buffer(); err != nil {
		return nil, fmt.Errorf("flush buffer: %w", err)
	}

	oggPages, err := s.oggPacker.ReadPages()
	if err != nil {
		return nil, fmt.Errorf("read pages: %w", err)
	}

	s.result = oggPages

	return s.result, nil
}

func (s *AudioBufferWriter) flushLastS16Buffer() error {
	defer func() {
		s.lastS16Buffer = s.lastS16Buffer[:0]
	}()

	if len(s.lastS16Buffer) == 0 {
		return nil
	}

	opusPackets, err := s.opusConverter.EncodeWithPadding(s.lastS16Buffer)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	for _, opusPacket := range opusPackets {
		if err := s.oggPacker.AddChunk(opusPacket, false, -1); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}

	return nil
}
