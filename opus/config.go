package opus

import "time"

const (
	SampleRate  = 48000
	frameSize   = time.Duration(60) * time.Millisecond
	numChannels = 1
	bufferSize  = 2048
)

type Config struct {
	SampleRate  int
	NumChannels int
	FrameSize   time.Duration
	BufferSize  int
}

func NewDefaultConfig() *Config {
	return &Config{
		SampleRate:  SampleRate,
		NumChannels: numChannels,
		FrameSize:   frameSize,
		BufferSize:  bufferSize,
	}
}
