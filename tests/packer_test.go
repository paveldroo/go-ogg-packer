package tests

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	packer "github.com/paveldroo/go-ogg-packer"
)

const (
	baseOggFilename = "testdata/base.ogg"
	rawOpusFilename = "testdata/48k_1ch_raw.opus"
)

func TestPacker1ch48khz(t *testing.T) {
	channelCount := 1
	sampleRate := 48000
	packer, err := packer.New(uint8(channelCount), uint32(sampleRate))
	if err != nil {
		t.Fatalf("create ogg packer: %s", err.Error())
	}

	rawOpusData, err := getRawOpusPackets(t)
	if err != nil {
		t.Fatalf("get result from audio buffer: %s", err.Error())
	}

	for _, packet := range rawOpusData {
		if err := packer.AddChunk(packet, false, -1); err != nil {
			t.Fatalf("send opus chunk to packer: %s", err.Error())
		}
	}

	oggData, err := packer.ReadPages()
	if err != nil {
		t.Fatalf("read all pages from packer: %s", err.Error())
	}

	fname := fmt.Sprintf("testdata/result/ogg_packer_result_%d.ogg", time.Now().UnixNano())
	mustWriteOggFile(fname, oggData)

	baseData, err := os.ReadFile(baseOggFilename)
	if err != nil {
		t.Fatalf("open base file: %s", err.Error())
	}

	if !reflect.DeepEqual(baseData, oggData) {
		t.Fatal("base data and test data are not equal")
	}
}

func getRawOpusPackets(t *testing.T) ([][]byte, error) {
	t.Helper()

	f, err := os.Open(rawOpusFilename)
	if err != nil {
		return nil, fmt.Errorf("read raw opus file: %w", err)
	}
	decoder := gob.NewDecoder(f)
	var audioData [][]byte
	if err := decoder.Decode(&audioData); err != nil {
		return nil, fmt.Errorf("decode data from file: %w", err)
	}

	return audioData, nil
}
