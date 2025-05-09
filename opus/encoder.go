package opus

import (
	"errors"
	"time"

	"gopkg.in/hraban/opus.v2"
)

var ErrTooLargeLastPacket = errors.New("last packet length is greater than frame size")

const (
	SampleRate  = 48000
	numChannels = 1
	frameSize   = time.Duration(60) * time.Millisecond
)

type Config struct {
	SampleRate  int
	NumChannels int
	FrameSize   time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		SampleRate:  SampleRate,
		NumChannels: numChannels,
		FrameSize:   frameSize,
	}
}

type Encoder struct {
	config           Config
	encoder          *encoderWrapper
	frameSizeSamples int
}

func NewEncoder(config Config) (*Encoder, error) {
	encoder, err := newEncoderWrapper(config.SampleRate, config.NumChannels, opus.AppAudio)
	if err != nil {
		return nil, err
	}

	return &Encoder{
		encoder:          encoder,
		config:           config,
		frameSizeSamples: FrameSizeSamples(config),
	}, nil
}

func (e *Encoder) Encode(samples []int16) ([][]byte, int, error) {
	var encoded [][]byte
	pos := 0
	for ; pos+e.frameSizeSamples <= len(samples); pos += e.frameSizeSamples {
		oneOpusPacket, err := e.encodeOneChunk(samples[pos : pos+e.frameSizeSamples])
		if err != nil {
			return [][]byte{}, 0, err
		}
		encoded = append(encoded, oneOpusPacket)
	}
	return encoded, pos, nil
}

func (e *Encoder) EncodeWithPadding(samples []int16) ([][]byte, error) {
	encoded, pos, err := e.Encode(samples)
	if err != nil {
		return nil, err
	}
	if len(samples) > pos {
		if len(samples)-pos > e.frameSizeSamples {
			return nil, ErrTooLargeLastPacket
		}
		samples = append(samples, make([]int16, e.frameSizeSamples-(len(samples)-pos))...)
		oneOpusPacket, err := e.encodeOneChunk(samples[pos : pos+e.frameSizeSamples])
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, oneOpusPacket)
	}
	return encoded, nil
}

func (e *Encoder) encodeOneChunk(samplesChunk []int16) ([]byte, error) {
	if len(samplesChunk) < e.frameSizeSamples {
		return []byte{}, nil
	}
	bufferSize := e.frameSizeSamples * 4
	oneOpusPacket := make([]byte, bufferSize)
	n, err := e.encoder.encode(samplesChunk[:e.frameSizeSamples], oneOpusPacket)
	if err != nil {
		return nil, err
	}
	oneOpusPacket = oneOpusPacket[:n]
	return oneOpusPacket, nil
}

func FrameSizeSamples(cfg Config) int {
	frameSizeMillis := cfg.FrameSize.Milliseconds()
	frameSizeSamples := float32(int64(cfg.NumChannels*cfg.SampleRate)*frameSizeMillis) / 1000
	return int(frameSizeSamples)
}
