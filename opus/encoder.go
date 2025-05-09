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

func NewDefaultConfig() *Config {
	return &Config{
		SampleRate:  SampleRate,
		NumChannels: numChannels,
		FrameSize:   frameSize,
	}
}

type Encoder struct {
	config           *Config
	encoder          *encoderWrapper
	frameSizeSamples int
}

func NewEncoder(config *Config) (*Encoder, error) {
	encoder, err := newEncoderWrapper(config.SampleRate, config.NumChannels, opus.AppAudio)
	if err != nil {
		return nil, err
	}

	frameSizeMillis := config.FrameSize.Milliseconds()
	frameSizeSamples := float32(int64(config.NumChannels*config.SampleRate)*frameSizeMillis) / 1000

	return &Encoder{
		encoder:          encoder,
		config:           config,
		frameSizeSamples: int(frameSizeSamples),
	}, nil
}

func (converter *Encoder) Encode(samples []int16) ([][]byte, int, error) {
	var encoded [][]byte
	pos := 0
	for ; pos+converter.frameSizeSamples <= len(samples); pos += converter.frameSizeSamples {
		oneOpusPacket, err := converter.encodeOneChunk(samples[pos : pos+converter.frameSizeSamples])
		if err != nil {
			return [][]byte{}, 0, err
		}
		encoded = append(encoded, oneOpusPacket)
	}
	return encoded, pos, nil
}

func (converter *Encoder) EncodeWithPadding(samples []int16) ([][]byte, error) {
	encoded, pos, err := converter.Encode(samples)
	if err != nil {
		return nil, err
	}
	if len(samples) > pos {
		if len(samples)-pos > converter.frameSizeSamples {
			return nil, ErrTooLargeLastPacket
		}
		samples = append(samples, make([]int16, converter.frameSizeSamples-(len(samples)-pos))...)
		oneOpusPacket, err := converter.encodeOneChunk(samples[pos : pos+converter.frameSizeSamples])
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, oneOpusPacket)
	}
	return encoded, nil
}

func (converter *Encoder) encodeOneChunk(samplesChunk []int16) ([]byte, error) {
	if len(samplesChunk) < converter.frameSizeSamples {
		return []byte{}, nil
	}
	bufferSize := converter.frameSizeSamples * 4
	oneOpusPacket := make([]byte, bufferSize)
	n, err := converter.encoder.encode(samplesChunk[:converter.frameSizeSamples], oneOpusPacket)
	if err != nil {
		return nil, err
	}
	oneOpusPacket = oneOpusPacket[:n]
	return oneOpusPacket, nil
}
