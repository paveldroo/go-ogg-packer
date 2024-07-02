package oggpacker

import "github.com/pion/opus"

type OggStreamState struct {
	bodyData       []byte // bytes from packet bodies
	bodyStorage    int    // storage elements allocated
	bodyFill       int    // elements stored; fill mark
	bodyReturned   int    // elements of fill returned
	lacingVals     int    // the values that will go to the segment table
	granuleVals    int64  // granulepos values for headers. Not compact this way, but it is simple coupled to the lacing fifo
	lacingStorage  int
	lacingFill     int
	lacingPacket   int
	lacingReturned int
	header         [282]byte // working space for header encode
	headerFill     int
	eos            int // set when we have buffered the last packet in the logical bitstream
	bos            int // set after we've written the initial page of a logical bitstream
	serialNo       int
	pageNo         int
	packetNo       int64 // sequence number for decode; the framing knows where there's a hole in the data, but we need coupling so that the codec (which is in a separate abstraction layer) also knows about the gap
	granulePos     int64
}

type Buffer struct {
	data      []byte
	len       uintptr
	readIndex uintptr
	alloc     uintptr
}

type Packer struct {
	channelCount uint8
	sampleRate   uint32
	packetNo     int64
	granulePos   int64
	streamState  *OggStreamState
	opusDecoder  *opus.Decoder
}

func NewPacker(sampleRate, numChannels int) (*Packer, error) {
	return nil, nil
}
