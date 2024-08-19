package cgo_oggpacker_test

import (
	"testing"

	"github.com/paveldroo/go-ogg-packer/lib"
	"github.com/paveldroo/go-ogg-packer/lib/cgo_oggpacker"
	"github.com/paveldroo/go-ogg-packer/lib/cgo_oggpacker/testdata"
)

func TestPacker_ReadAudioData(t *testing.T) {
	sampleRate := 8000
	numChannels := 1

	oggPacker, err := cgo_oggpacker.New(sampleRate, numChannels)
	if err != nil {
		t.Fatalf("create oggPacker: %s", err.Error())
	}
	chunkSender := testdata.AudioByChunks()

	var resData []byte

	for _, chunk := range chunkSender {
		if err := oggPacker.AddChunk(chunk); err != nil {
			t.Fatalf("add chunk: %s", err.Error())
		}
	}

	resData, err = oggPacker.ReadAudioData()
	if err != nil {
		t.Fatalf("readAudioData from packer: %s", err.Error())
	}

	lib.MustWriteResultFile(resData)

	refData := testdata.RefOGGData()

	if !testdata.CompareOggAudio(resData, refData) {
		t.Fatalf("result and reference audio files are not the same")
	}
}
