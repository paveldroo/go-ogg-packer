package opus_tools

import (
	"fmt"

	"gopkg.in/hraban/opus.v2"
)

type Converter struct {
	config           *Config
	encoder          *opus.Encoder
	frameSizeSamples int
}

func NewOpusConverter(config *Config) (*Converter, error) {
	encoder, err := opus.NewEncoder(config.SampleRate, config.NumChannels, opus.AppAudio)
	if err != nil {
		return nil, fmt.Errorf("create encoder: %w", err)
	}

	frameSizeMillis := config.FrameSize.Milliseconds()
	frameSizeSamples := float32(int64(config.NumChannels*config.SampleRate)*frameSizeMillis) / 1000

	return &Converter{
		encoder:          encoder,
		config:           config,
		frameSizeSamples: int(frameSizeSamples),
	}, nil
}

func (converter *Converter) EncodeOneChunk(samplesChunk []int16) ([]byte, error) {
	if len(samplesChunk) < converter.frameSizeSamples {
		return []byte{}, nil
	}
	bufferSize := converter.frameSizeSamples * 4
	oneOpusPacket := make([]byte, bufferSize)
	n, err := converter.encoder.Encode(samplesChunk[:converter.frameSizeSamples], oneOpusPacket)
	if err != nil {
		return nil, err
	}
	oneOpusPacket = oneOpusPacket[:n]
	return oneOpusPacket, nil
}

func (converter *Converter) Encode(samples []int16) ([][]byte, int, error) {
	var encoded [][]byte
	pos := 0
	for ; pos+converter.frameSizeSamples <= len(samples); pos += converter.frameSizeSamples {
		oneOpusPacket, err := converter.EncodeOneChunk(samples[pos : pos+converter.frameSizeSamples])
		if err != nil {
			return [][]byte{}, 0, err
		}
		encoded = append(encoded, oneOpusPacket)
	}
	return encoded, pos, nil
}

func (converter *Converter) EncodeWithPadding(samples []int16) ([][]byte, error) {
	encoded, pos, err := converter.Encode(samples)
	if err != nil {
		return nil, err
	}
	if len(samples) > pos {
		if len(samples)-pos > converter.frameSizeSamples {
			return nil, ErrTooLargeLastPacket
		}
		samples = append(samples, make([]int16, converter.frameSizeSamples-(len(samples)-pos))...)
		oneOpusPacket, err := converter.EncodeOneChunk(samples[pos : pos+converter.frameSizeSamples])
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, oneOpusPacket)
	}
	return encoded, nil
}
