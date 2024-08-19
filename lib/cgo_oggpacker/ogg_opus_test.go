package cgo_oggpacker_test

import (
	"testing"

	"github.com/paveldroo/go-ogg-packer/lib/cgo_oggpacker"
	"github.com/paveldroo/go-ogg-packer/testdata"
)

func TestPacker_ReadPages(t *testing.T) {
	sampleRate := 8000
	numChannels := 1

	oggPacker, err := cgo_oggpacker.New(sampleRate, numChannels)
	if err != nil {
		t.Fatalf("create oggPacker: %s", err)
	}
	chunkSender := testdata.NewChunkSender(testdata.MustReferencePath(), 512)

	for chunk := range chunkSender {
		oggPacker.AddChunk(chunk)
		oggPacker.AddChunk(chunk)
		oggPacker.AddChunk(chunk)
		oggPacker.AddChunk(chunk)

		// TODO: need to add chunkWrapper with more convenient API for adding chunks, CGO implementation is too low level for testing.
	}
}
