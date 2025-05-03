package tests

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
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

	rawOpusData := getRawOpusPackets(t)
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
	writeOggFile(t, fname, oggData)

	baseData, err := os.ReadFile(baseOggFilename)
	if err != nil {
		t.Fatalf("open base file: %s", err.Error())
	}

	if !reflect.DeepEqual(baseData, oggData) {
		t.Fatal("base data and test data are not equal")
	}
}

func getRawOpusPackets(t *testing.T) [][]byte {
	t.Helper()

	f, err := os.Open(rawOpusFilename)
	if err != nil {
		t.Fatalf("read raw opus file: %s", err.Error())
	}
	decoder := gob.NewDecoder(f)
	var audioData [][]byte
	if err := decoder.Decode(&audioData); err != nil {
		t.Fatalf("decode data from file: %s", err.Error())
	}

	return audioData
}

func writeOggFile(t *testing.T, name string, data []byte) {
	t.Helper()
	wDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, name)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		t.Fatalf("write result file: %s", err.Error())
	}
}
