package cgo_oggpacker

import (
	"testing"

	"github.com/paveldroo/go-ogg-packer/opus"
	"gitlab.tcsbank.ru/speech/libanysound/ogg_packer"
)

// func TestOGGPackerWrapper(t *testing.T) {
// 	wrapper, err := newOggPackerWrapper(sampleRate, 1)
// 	if err != nil {
// 		t.Fatalf("create ogg packer wrapper: %s", err.Error())
// 	}

// 	converter, err := opus.NewOpusConverter(opus.NewDefaultConfig())
// 	if err != nil {
// 		t.Fatalf("create opus converter: %s", err.Error())
// 	}

// 	s16 := S16FromWav()

// 	currentOpusPackets, _, err := converter.Encode(s16)
// 	if err != nil {
// 		t.Fatalf("encode s16 to opus: %s", err.Error())
// 	}

// 	for _, opusPacket := range currentOpusPackets {
// 		if err := wrapper.addChunk(opusPacket); err != nil {
// 			t.Fatalf("send opusPacket to addChunk: %s", err.Error())
// 		}
// 	}

// 	oggPages, err := wrapper.packer.ReadPages()
// 	if err != nil {
// 		t.Fatalf("read pages from ogg packer: %s", err.Error())
// 	}

// 	flushedPages, err := wrapper.packer.FlushPages()
// 	if err != nil {
// 		t.Fatalf("flush pages from ogg packer: %s", err.Error())
// 	}

// 	oggPages = append(oggPages, flushedPages...)

// 	// fmt.Println("res", res)

// 	mustWriteOpusFile("testdata/ogg_packer_result.opus", oggPages)
// }

func TestOGGPacker(t *testing.T) {
	packer, err := ogg_packer.New(sampleRate, 1)
	if err != nil {
		t.Fatalf("create ogg packer wrapper: %s", err.Error())
	}

	converter, err := opus.NewOpusConverter(opus.NewDefaultConfig())
	if err != nil {
		t.Fatalf("create opus converter: %s", err.Error())
	}

	s16 := S16FromWav()

	currentOpusPackets, pos, err := converter.Encode(s16)
	if err != nil {
		t.Fatalf("encode s16 to opus: %s", err.Error())
	}

	for _, opusPacket := range currentOpusPackets {
		if err := packer.AddChunk(opusPacket, false, pos); err != nil {
			t.Fatalf("send opusPacket to addChunk: %s", err.Error())
		}
	}

	oggPages, err := packer.ReadPages()
	if err != nil {
		t.Fatalf("read pages from ogg packer: %s", err.Error())
	}

	flushedPages, err := packer.FlushPages()
	if err != nil {
		t.Fatalf("flush pages from ogg packer: %s", err.Error())
	}

	oggPages = append(oggPages, flushedPages...)

	// fmt.Println("res", res)

	mustWriteOpusFile("testdata/ogg_packer_result.opus", oggPages)
}
