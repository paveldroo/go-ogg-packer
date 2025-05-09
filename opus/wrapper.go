package opus

import (
	"fmt"
	"sync"

	"gopkg.in/hraban/opus.v2"
)

// newEncoderWrapper creates concurrent safe Opus encoder
func newEncoderWrapper(sampleRate, channels int, application opus.Application) (*encoderWrapper, error) {
	encoder, err := opus.NewEncoder(sampleRate, channels, application)
	if err != nil {
		return nil, fmt.Errorf("create encoder: %w", err)
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

	val, err := s.encoder.Encode(pcm, data)
	if err != nil {
		return 0, fmt.Errorf("encode: %w", err)
	}

	return val, nil
}
