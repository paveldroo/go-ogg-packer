package opus_tools

import "time"

const (
	SampleRate  = 48000
	numChannels = 1
	frameSize   = time.Duration(60) * time.Millisecond
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
