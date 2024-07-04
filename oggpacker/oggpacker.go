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

type Packer struct {
	channelCount uint8
	sampleRate   uint32
	packetNo     int64
	granulePos   int64
	streamState  *OggStreamState
	opusDecoder  *opus.Decoder
}

func NewPacker(sampleRate, numChannels int) (*Packer, error) {
	d := opus.NewDecoder()
	ss := OggStreamState{
		bodyData:       nil,
		bodyStorage:    0,
		bodyFill:       0,
		bodyReturned:   0,
		lacingVals:     0,
		granuleVals:    0,
		lacingStorage:  0,
		lacingFill:     0,
		lacingPacket:   0,
		lacingReturned: 0,
		header:         [282]byte{},
		headerFill:     0,
		eos:            0,
		bos:            0,
		serialNo:       0,
		pageNo:         0,
		packetNo:       0,
		granulePos:     0,
	}
	return &Packer{
		channelCount: uint8(numChannels),
		sampleRate:   uint32(sampleRate),
		packetNo:     0,
		granulePos:   0,
		streamState:  &ss,
		opusDecoder:  &d,
	}, nil
}
