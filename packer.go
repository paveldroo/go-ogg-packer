package packer

import (
	"fmt"
	"slices"
	"time"

	"github.com/paveldroo/go-ogg-packer/ogg"
	"github.com/paveldroo/go-ogg-packer/opus"
)

type Packer struct {
	result      []byte
	opusEncoder *opus.Encoder
	oggPacker   *ogg.Packer
	pcmBuffer   []int16
}

func New() (*Packer, error) {
	cfg := opus.NewDefaultConfig()
	return newPacker(cfg)
}

func NewWithConfig(sampleRate int, numChannels int, frameSize time.Duration) (*Packer, error) {
	if !slices.Contains(opus.SupportedSampleRate, sampleRate) {
		return nil, fmt.Errorf("sample rate must be one of the supported sample rates (8000, 12000, 16000, 24000, 48000)")
	}
	if !slices.Contains(opus.SupportedNumChannels, numChannels) {
		return nil, fmt.Errorf("num channels must be one of the supported num channels (1,2)")
	}
	if !slices.Contains(opus.SupportedFrameSize, float64(frameSize)/float64(time.Millisecond)) {
		return nil, fmt.Errorf("framesize must be (2.5, 5, 10, 20, 40 or 60) ms")
	}
	return newPacker(opus.Config{
		SampleRate:  sampleRate,
		NumChannels: numChannels,
		FrameSize:   frameSize,
	})
}
func newPacker(cfg opus.Config) (*Packer, error) {
	encoder, err := opus.NewEncoder(cfg)
	if err != nil {
		return nil, fmt.Errorf("create opus encoder: %s", err)
	}

	packer, err := ogg.New(uint8(cfg.NumChannels), uint32(cfg.SampleRate))
	if err != nil {
		return nil, fmt.Errorf("create ogg packer: %w", err)
	}

	return &Packer{
		opusEncoder: encoder,
		oggPacker:   packer,
	}, nil
}

func (s *Packer) SendPCMChunk(chunk []int16) error {
	s.pcmBuffer = append(s.pcmBuffer, chunk...)
	currentOpusPackets, pos, err := s.opusEncoder.Encode(s.pcmBuffer)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	s.pcmBuffer = s.pcmBuffer[pos:]
	for _, opusPacket := range currentOpusPackets {
		if err := s.oggPacker.AddChunk(opusPacket, false, pos); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}
	return nil
}

func (s *Packer) GetResult() ([]byte, error) {
	defer s.oggPacker.Close()

	if err := s.flushPCMBuffer(); err != nil {
		return nil, fmt.Errorf("flush buffer: %w", err)
	}

	oggPages, err := s.oggPacker.ReadPages()
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

	for _, opusPacket := range opusPackets {
		if err := s.oggPacker.AddChunk(opusPacket, false, -1); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}

	return nil
}
