package opus

import "time"

const (
	frameSize   = time.Duration(60) * time.Millisecond
	sampleRate  = 48000
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
		SampleRate:  sampleRate,
		NumChannels: numChannels,
		FrameSize:   frameSize,
		BufferSize:  bufferSize,
	}
}
