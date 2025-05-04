package tests

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	packer "github.com/paveldroo/go-ogg-packer"
)

func TestPacker(t *testing.T) {
	tests := []struct {
		name       string
		channels   int
		sampleRate int
	}{
		{
			name:       "48k 1ch",
			channels:   1,
			sampleRate: 48000,
		},
		// {
		// 	name:       "48k 2ch",
		// 	channels:   2,
		// 	sampleRate: 48000,
		// },
		// {
		// 	name:       "8k 1ch",
		// 	channels:   1,
		// 	sampleRate: 8000,
		// },
		// {
		// 	name:       "8k 2ch",
		// 	channels:   2,
		// 	sampleRate: 8000,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseFilename := fmt.Sprintf("%dk_%dch", tt.sampleRate, tt.channels)

			packer, err := packer.New(uint8(tt.channels), uint32(tt.sampleRate))
			if err != nil {
				t.Fatalf("create ogg packer: %s", err.Error())
			}

			opusFilename := fmt.Sprintf("testdata/opus_raw/%s.opus_raw", baseFilename)
			rawOpusData := getRawOpusPackets(t, opusFilename)
			for _, packet := range rawOpusData {
				if err := packer.AddChunk(packet, false, -1); err != nil {
					t.Fatalf("send opus chunk to packer: %s", err.Error())
				}
			}

			oggData, err := packer.ReadPages()
			if err != nil {
				t.Fatalf("read all pages from packer: %s", err.Error())
			}

			fname := fmt.Sprintf("testdata/%s.ogg", baseFilename)
			writeOggFile(t, fname, oggData)

			refFilename := fmt.Sprintf("testdata/want/%s.ogg", baseFilename)
			refData, err := os.ReadFile(refFilename)
			if err != nil {
				t.Fatalf("open reference file: %s", err.Error())
			}

			if !reflect.DeepEqual(refData, oggData) {
				t.Fatal("base data and test data are not equal")
			}
		})
	}
}

func getRawOpusPackets(t *testing.T, fname string) [][]byte {
	t.Helper()

	f, err := os.Open(fname)
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

func writeOggFile(t *testing.T, fname string, data []byte) {
	t.Helper()
	wDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, fname)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		t.Fatalf("write result file: %s", err.Error())
	}
}
