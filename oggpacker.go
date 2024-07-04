package go_ogg_packer

import (
	"github.com/pion/opus"
	"math/rand"
	"time"
)

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
	buffer       *Buffer
	opusDecoder  *opus.Decoder
}

func NewPacker(sampleRate, numChannels int) (*Packer, error) {
	d := opus.NewDecoder()
	sn := rand.New(rand.NewSource(time.Now().UTC().Unix() % 0x80000000)).Int()
	b := Buffer{
		data:      nil,
		len:       0,
		readIndex: 0,
		alloc:     0,
	}
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
		serialNo:       sn,
		pageNo:         0,
		packetNo:       0,
		granulePos:     0,
	}

	p := Packer{
		channelCount: 0,
		sampleRate:   0,
		packetNo:     0,
		granulePos:   0,
		streamState:  nil,
		buffer:       nil,
		opusDecoder:  nil,
	}

	status := p.Init(sampleRate, numChannels, sn)

	return &p, nil
}

func (p *Packer) Init(sampleRate, numChannels, serialNo int) int {
	p.channelCount = uint8(numChannels)
	p.sampleRate = uint32(sampleRate)
	p.packetNo = 1
	p.granulePos = 0

}
