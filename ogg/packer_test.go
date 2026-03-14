package ogg_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/paveldroo/go-ogg-packer/ogg"
	"github.com/paveldroo/go-ogg-packer/opus"
	"github.com/paveldroo/go-ogg-packer/tests/testutil"
)

// TestPacker compares OGG packer output against reference OGG files.
// References must be generated before running: task generate-ref
// Or run everything together: task test
func TestPacker(t *testing.T) {
	tests := []struct {
		name       string
		fileBase   string
		channels   int
		sampleRate int
	}{
		{
			name:       "8k 1ch",
			fileBase:   "8k_1ch",
			channels:   1,
			sampleRate: 8000,
		},
		{
			name:       "12k 1ch",
			fileBase:   "12k_1ch",
			channels:   1,
			sampleRate: 12000,
		},
		{
			name:       "16k 1ch",
			fileBase:   "16k_1ch",
			channels:   1,
			sampleRate: 16000,
		},
		{
			name:       "24k 1ch",
			fileBase:   "24k_1ch",
			channels:   1,
			sampleRate: 24000,
		},
		{
			name:       "48k 1ch",
			fileBase:   "48k_1ch",
			channels:   1,
			sampleRate: 48000,
		},
		{
			name:       "8k 2ch",
			fileBase:   "8k_2ch",
			channels:   2,
			sampleRate: 8000,
		},
		{
			name:       "12k 2ch",
			fileBase:   "12k_2ch",
			channels:   2,
			sampleRate: 12000,
		},
		{
			name:       "16k 2ch",
			fileBase:   "16k_2ch",
			channels:   2,
			sampleRate: 16000,
		},
		{
			name:       "24k 2ch",
			fileBase:   "24k_2ch",
			channels:   2,
			sampleRate: 24000,
		},
		{
			name:       "48k 2ch",
			fileBase:   "48k_2ch",
			channels:   2,
			sampleRate: 48000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packer, err := ogg.New(uint8(tt.channels), uint32(tt.sampleRate))
			if err != nil {
				t.Fatalf("create ogg packer: %s", err.Error())
			}

			cfg := opus.Config{
				SampleRate:  tt.sampleRate,
				NumChannels: tt.channels,
				FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
			}
			frameSizeSamples := opus.FrameSizeSamples(cfg)

			opusFilename := fmt.Sprintf("../tests/testdata/opus_raw/%s.opus_raw", tt.fileBase)
			rawOpusData := testutil.RawOpusPackets(t, opusFilename)
			for _, packet := range rawOpusData {
				if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
					t.Fatalf("send opus chunk to packer: %s", err.Error())
				}
			}

			oggData, err := packer.ReadPages()
			if err != nil {
				t.Fatalf("read all pages from packer: %s", err.Error())
			}

			refFilename := fmt.Sprintf("../tests/testdata/want/ogg/%s.ogg", tt.fileBase)
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

func TestPacker_EOS(t *testing.T) {
	t.Run("ReadPages patches EOS on last page", func(t *testing.T) {
		packer, err := ogg.New(1, 48000)
		if err != nil {
			t.Fatalf("create ogg packer: %s", err)
		}

		opusFilename := "../tests/testdata/opus_raw/48k_1ch.opus_raw"
		packets := testutil.RawOpusPackets(t, opusFilename)
		for _, packet := range packets {
			if err := packer.AddChunk(packet, false, 2880); err != nil {
				t.Fatalf("add chunk: %s", err)
			}
		}

		oggData, err := packer.ReadPages()
		if err != nil {
			t.Fatalf("read pages: %s", err)
		}

		assertLastPageHasEOS(t, oggData)
	})

	t.Run("explicit EOS via AddChunk", func(t *testing.T) {
		packer, err := ogg.New(1, 48000)
		if err != nil {
			t.Fatalf("create ogg packer: %s", err)
		}

		opusFilename := "../tests/testdata/opus_raw/48k_1ch.opus_raw"
		packets := testutil.RawOpusPackets(t, opusFilename)
		for i, packet := range packets {
			eos := i == len(packets)-1
			if err := packer.AddChunk(packet, eos, 2880); err != nil {
				t.Fatalf("add chunk: %s", err)
			}
		}

		oggData, err := packer.ReadPages()
		if err != nil {
			t.Fatalf("read pages: %s", err)
		}

		assertLastPageHasEOS(t, oggData)
	})
}

func assertLastPageHasEOS(t *testing.T, data []byte) {
	t.Helper()

	lastPageStart := -1
	for i := len(data) - 4; i >= 0; i-- {
		if string(data[i:i+4]) == "OggS" {
			lastPageStart = i
			break
		}
	}
	if lastPageStart < 0 {
		t.Fatal("no OggS page found in output")
	}
	if data[lastPageStart+5]&0x04 == 0 { // 0x04 = EOS flag per RFC 3533
		t.Fatal("last page does not have EOS flag set")
	}
}
