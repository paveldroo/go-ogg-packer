package opus

import (
	"errors"
	"fmt"

	"gopkg.in/hraban/opus.v2"
)

// Caution: Do not use opus.Converter with multiple audio streams to avoid sound artefacts.

type Converter struct {
	config           *Config
	encoder          *encoderWrapper
	decoder          *decoderWrapper
	frameSizeSamples int
}

func NewOpusConverter(config *Config) (*Converter, error) {
	// TODO Consider using opus.AppVoIP instead of opus.AppAudio
	encoder, err := newEncoderWrapper(config.SampleRate, config.NumChannels, opus.AppAudio)
	if err != nil {
		return nil, err
	}
	decoder, err := newDecoderWrapper(config.SampleRate, config.NumChannels)
	if err != nil {
		return nil, err
	}
	frameSizeSamples := int(float32(int64(config.NumChannels*config.SampleRate)*config.FrameSize.Milliseconds()) / 1000)
	return &Converter{
		encoder:          encoder,
		decoder:          decoder,
		config:           config,
		frameSizeSamples: frameSizeSamples,
	}, nil
}

func (converter *Converter) EncodeOneChunk(samplesChunk []int16) ([]byte, error) {
	if len(samplesChunk) < converter.frameSizeSamples {
		return []byte{}, nil
	}
	oneOpusPacket := make([]byte, converter.config.BufferSize)
	n, err := converter.encoder.encode(samplesChunk[:converter.frameSizeSamples], oneOpusPacket)
	if err != nil {
		return nil, err
	}
	oneOpusPacket = oneOpusPacket[:n]
	return oneOpusPacket, nil
}

func (converter *Converter) Encode(samples []int16) ([][]byte, int, error) {
	var encoded [][]byte
	pos := 0
	for ; pos+converter.frameSizeSamples < len(samples); pos += converter.frameSizeSamples {
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
			return nil, fmt.Errorf("last packet length is greater than frame size")
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

func (converter *Converter) DecodeOneChunk(encodedPacket []byte) ([]int16, error) {
	if len(encodedPacket) > converter.config.BufferSize {
		return []int16{}, errors.New(
			fmt.Sprintf("Opus packet length should not be greater than %d bytes.", converter.config.BufferSize))
	}
	pcm := make([]int16, converter.frameSizeSamples)
	n, err := converter.decoder.decode(encodedPacket, pcm)
	if err != nil {
		return nil, err
	}
	pcm = pcm[:n]
	return pcm, nil
}

func (converter *Converter) Decode(encoded [][]byte) ([]int16, error) {
	samples := make([]int16, 0, int(float64(len(encoded))*float64(converter.config.BufferSize)*converter.config.FrameSize.Seconds()))
	for i := 0; i < len(encoded); i++ {
		pcm, err := converter.DecodeOneChunk(encoded[i])
		if err != nil {
			return []int16{}, err
		}
		samples = append(samples, pcm...)
	}
	return samples, nil
}
