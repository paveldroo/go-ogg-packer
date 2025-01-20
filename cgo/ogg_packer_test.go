package cgo_oggpacker

import (
	"fmt"
	"testing"
)

func Test_OGGPacker(t *testing.T) {
	wrapper, err := newOggPackerWrapper(sampleRate, 1)
	if err != nil {
		t.Fatalf("create ogg packer wrapper: %s", err.Error())
	}

	d := AudioByChunks()

	for _, chunk := range d {
		wrapper.addChunk(chunk)
	}

	res, err := wrapper.readAudioData()
	if err != nil {
		t.Fatalf("read audio data: %s", err.Error())
	}

	fmt.Println("res", res)

	mustWriteOpusFile(res)
}
