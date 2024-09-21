package audio_buffer_writer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/paveldroo/go-ogg-packer/lib/audio_buffer_writer/opus"
	"github.com/paveldroo/go-ogg-packer/lib/audio_buffer_writer/utils"
	"github.com/paveldroo/go-ogg-packer/lib/cgo_oggpacker"
)

const (
	Linear16 = "linear16"
	Alaw     = "alaw"
	RawOpus  = "opus"
)

func NewBufferWriter(opusConfig *opus.Config) (*AudioBufferWriter, error) {
	opusConverter, err := opus.NewOpusConverter(opusConfig)
	if err != nil {
		return nil, err
	}
	oggPacker, err := cgo_oggpacker.New(opusConfig.SampleRate, opusConfig.NumChannels)
	if err != nil {
		return nil, err
	}
	return &AudioBufferWriter{
		AudioEncoding:     RawOpus,
		opusConverter:     opusConverter,
		oggPacker:         oggPacker,
		lastS16Buffer:     make([]int16, 0),
		mutex:             new(sync.Mutex),
		sentTotalDuration: &totalDuration{},
	}, nil
}

type totalDuration struct {
	totalDurationCached    float64
	totalDurationNotCached float64
}

type AudioBufferWriter struct {
	AudioEncoding     string
	result            []byte
	pcmBuffer         bytes.Buffer
	opusConverter     *opus.Converter
	oggPacker         *cgo_oggpacker.Packer
	lastS16Buffer     []int16
	finalized         bool
	mutex             *sync.Mutex
	sentTotalDuration *totalDuration
}

func (s *AudioBufferWriter) sendS16Chunk(chunk []int16, cacheHit bool) error {
	if s.finalized {
		return fmt.Errorf("AudioBufferWriter already finalized")
	}
	switch encoding := s.AudioEncoding; encoding {
	case Linear16:
		if binary.Write(&s.pcmBuffer, binary.LittleEndian, chunk) != nil {
			return errors.New("error writing to buffer")
		}
	case Alaw:
		if binary.Write(&s.pcmBuffer, binary.LittleEndian, utils.S16ToALaw(chunk)) != nil {
			return errors.New("error writing to buffer")
		}
	case RawOpus:
		s.lastS16Buffer = append(s.lastS16Buffer, chunk...)
		currentOpusPackets, pos, err := s.opusConverter.Encode(s.lastS16Buffer)
		if err != nil {
			return err
		}
		s.lastS16Buffer = s.lastS16Buffer[pos:]
		for _, opusPacket := range currentOpusPackets {
			if err := s.oggPacker.AddChunk(opusPacket, false, pos); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("encoding is not supported: %s", encoding)
	}
	return nil
}

func (s *AudioBufferWriter) getResult() ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.finalized {
		defer func() {
			s.oggPacker.Close()
			s.finalized = true
		}()
		if err := s.flushLastS16Buffer(); err != nil {
			return nil, err
		}
		switch s.AudioEncoding {
		case Linear16, Alaw:
			s.result = s.pcmBuffer.Bytes()
		case RawOpus:
			oggPages, err := s.oggPacker.ReadPages()
			if err != nil {
				return nil, err
			}
			flushedOggPages, err := s.oggPacker.FlushPages()
			if err != nil {
				return nil, err
			}
			oggPages = append(oggPages, flushedOggPages...)
			s.result = oggPages
		default:
			return nil, fmt.Errorf("encoding is not supported")
		}
	}
	return s.result, nil
}

func (s *AudioBufferWriter) flushLastS16Buffer() error {
	if s.finalized {
		return fmt.Errorf("AudioBufferWriter already finalized")
	}
	defer func() {
		s.lastS16Buffer = s.lastS16Buffer[:0]
	}()
	if len(s.lastS16Buffer) > 0 {
		if s.AudioEncoding == Linear16 {
			if binary.Write(&s.pcmBuffer, binary.LittleEndian, s.lastS16Buffer) != nil {
				return errors.New("error writing to buffer")
			}
		} else {
			opusPackets, err := s.opusConverter.EncodeWithPadding(s.lastS16Buffer)
			if err != nil {
				return err
			}
			for _, opusPacket := range opusPackets {
				if err := s.oggPacker.AddChunk(opusPacket, false, -1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
