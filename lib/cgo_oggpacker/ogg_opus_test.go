package cgo_oggpacker_test

import (
	"fmt"
	"testing"

	"github.com/paveldroo/go-ogg-packer/lib"
	"github.com/paveldroo/go-ogg-packer/lib/cgo_oggpacker"
	"github.com/paveldroo/go-ogg-packer/lib/cgo_oggpacker/testdata"
)

func TestPacker_ReadAudioData(t *testing.T) {
	sampleRate := 48000
	numChannels := 1

	oggPacker, err := cgo_oggpacker.New(sampleRate, numChannels)
	if err != nil {
		t.Fatalf("create oggPacker: %s", err.Error())
	}
	chunkSender := testdata.AudioByChunks()

	var resData []byte

	for i, chunk := range chunkSender {
		fmt.Printf("Processing chunk %d, size: %d\n", i, len(chunk))
		if err := oggPacker.AddChunk(chunk); err != nil {
			t.Fatalf("add chunk: %s", err.Error())
		}
	}

	resData, err = oggPacker.ReadAudioData()
	if err != nil {
		t.Fatalf("readAudioData from packer: %s", err.Error())
	}

	fmt.Printf("Encoded OGG data size: %d bytes\n", len(resData))

	lib.MustWriteResultFile(resData)

	//refData := testdata.RefOGGData()
	//
	//if !testdata.CompareOggAudio(resData, refData) {
	//	t.Fatalf("result and reference audio files are not the same")
	//}
}

//func TestConvertBytes(t *testing.T) {
//	chunkSender := testdata.AudioByChunks()
//
//	var resData []byte
//
//	for _, chunk := range chunkSender {
//		resData = append(resData, chunk...)
//	}
//
//	lib.MustWriteResultFile(resData)
//
//}
