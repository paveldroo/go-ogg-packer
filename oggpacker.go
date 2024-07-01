package oggpacker

import "github.com/pion/opus"

type OggStreamState struct {
	// ogg_stream_state stream_state;
}

type Buffer struct {
	data      []byte
	len       uintptr
	readIndex uintptr
	alloc     uintptr
}

type OggPacker struct {
	channelCount uint8
	sampleRate   uint32
	packetNo     int64
	granulePos   int64
	streamState  OggStreamState
	opusDecoder  *opus.Decoder
}

func NewOggPacker(sampleRate, numChannels int) (*OggPacker, error) {
	return nil, nil
}
