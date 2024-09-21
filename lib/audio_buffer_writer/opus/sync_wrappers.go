package opus

import (
	"sync"

	"gopkg.in/hraban/opus.v2"
)

// Opus encoder and decoder can produce undefined behavior when used concurrently.
// These wrappers are defined to avoid it.

func newEncoderWrapper(sampleRate, channels int, application opus.Application) (*encoderWrapper, error) {
	encoder, err := opus.NewEncoder(sampleRate, channels, application)
	if err != nil {
		return nil, err
	}
	return &encoderWrapper{
		encoder: encoder,
		mutex:   new(sync.Mutex),
	}, nil
}

type encoderWrapper struct {
	encoder *opus.Encoder
	mutex   *sync.Mutex
}

func (s *encoderWrapper) encode(pcm []int16, data []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.encoder.Encode(pcm, data)
}

func newDecoderWrapper(sampleRate, channels int) (*decoderWrapper, error) {
	decoder, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
		return nil, err
	}
	return &decoderWrapper{
		decoder: decoder,
		mutex:   new(sync.Mutex),
	}, nil
}

type decoderWrapper struct {
	decoder *opus.Decoder
	mutex   *sync.Mutex
}

func (s *decoderWrapper) decode(data []byte, pcm []int16) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.decoder.Decode(data, pcm)
}
