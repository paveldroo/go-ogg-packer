package cgo_oggpacker

import (
	"fmt"

	"github.com/paveldroo/go-ogg-packer/opus"
	"gitlab.tcsbank.ru/speech/libanysound/ogg_packer"
)

type AudioBufferWriter struct {
	result        []byte
	opusConverter *opus.Converter
	oggPacker     *ogg_packer.Packer
	lastS16Buffer []int16
}

func NewAudioBuffer(
	opusConverter *opus.Converter,
	oggPacker *ogg_packer.Packer,
) *AudioBufferWriter {
	return &AudioBufferWriter{
		opusConverter: opusConverter,
		oggPacker:     oggPacker,
	}
}

func (s *AudioBufferWriter) sendS16Chunk(chunk []int16) error { //nolint:cyclop // SIMPLIFY
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

func (s *AudioBufferWriter) getResult() ([]byte, error) {
	defer s.oggPacker.Close()

	if err := s.flushLastS16Buffer(); err != nil {
		return nil, fmt.Errorf("flush buffer: %w", err)
	}

	oggPages, err := s.oggPacker.ReadPages()
	if err != nil {
		return nil, fmt.Errorf("read pages: %w", err)
	}
	flushedOggPages, err := s.oggPacker.FlushPages()
	if err != nil {
		return nil, fmt.Errorf("flush pages: %w", err)
	}
	oggPages = append(oggPages, flushedOggPages...)
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
