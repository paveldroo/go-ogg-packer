package ogg_test

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/paveldroo/go-ogg-packer/ogg"
	"github.com/paveldroo/go-ogg-packer/opus"
)

func TestPacker(t *testing.T) {
	genNewReference := os.Getenv("GENERATE_NEW_REFERENCE")

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

			opusFilename := fmt.Sprintf("testdata/opus_raw/%s.opus_raw", tt.fileBase)
			rawOpusData := rawOpusPackets(t, opusFilename)
			for _, packet := range rawOpusData {
				if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
					t.Fatalf("send opus chunk to packer: %s", err.Error())
				}
			}

			oggData, err := packer.ReadPages()
			if err != nil {
				t.Fatalf("read all pages from packer: %s", err.Error())
			}

			if genNewReference != "" {
				writeOggFile(t, fmt.Sprintf("testdata/want/%s.ogg", tt.fileBase), oggData)
				t.Logf("generated reference file: testdata/want/%s.ogg", tt.fileBase)
				return
			}

			refFilename := fmt.Sprintf("testdata/want/%s.ogg", tt.fileBase)
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

// TestGenerateOpusRaw generates .opus_raw fixture files by encoding PCM source data.
// Run with: GENERATE_OPUS_RAW=true go test ./ogg/ -run TestGenerateOpusRaw -v
func TestGenerateOpusRaw(t *testing.T) {
	if os.Getenv("GENERATE_OPUS_RAW") == "" {
		t.Skip("set GENERATE_OPUS_RAW=true to generate opus_raw fixtures")
	}

	rates := []struct {
		sampleRate int
		fileBase   string
		pcmSource  string
	}{
		{8000, "8k_1ch", "../testdata/8k_1ch.pcm"},
		{12000, "12k_1ch", "../testdata/12k_1ch.pcm"},
		{16000, "16k_1ch", "../testdata/16k_1ch.pcm"},
		{24000, "24k_1ch", "../testdata/24k_1ch.pcm"},
		{48000, "48k_1ch", "../testdata/48k_1ch.pcm"},
	}

	for _, r := range rates {
		t.Run(r.fileBase, func(t *testing.T) {
			cfg := opus.Config{
				SampleRate:  r.sampleRate,
				NumChannels: opus.NumChannels,
				FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
			}
			encoder, err := opus.NewEncoder(cfg)
			if err != nil {
				t.Fatalf("create encoder: %s", err)
			}

			pcmBytes, err := os.ReadFile(r.pcmSource)
			if err != nil {
				t.Fatalf("read pcm: %s", err)
			}

			samples := make([]int16, len(pcmBytes)/2)
			for i := range samples {
				samples[i] = int16(pcmBytes[2*i]) | int16(pcmBytes[2*i+1])<<8
			}

			packets, err := encoder.EncodeWithPadding(samples)
			if err != nil {
				t.Fatalf("encode: %s", err)
			}

			outPath := fmt.Sprintf("testdata/opus_raw/%s.opus_raw", r.fileBase)
			f, err := os.Create(outPath)
			if err != nil {
				t.Fatalf("create file: %s", err)
			}
			defer f.Close()

			if err := gob.NewEncoder(f).Encode(packets); err != nil {
				t.Fatalf("gob encode: %s", err)
			}

			t.Logf("generated %s (%d packets)", outPath, len(packets))
		})
	}
}

func TestPacker_EOS(t *testing.T) {
	t.Run("ReadPages patches EOS on last page", func(t *testing.T) {
		packer, err := ogg.New(1, 48000)
		if err != nil {
			t.Fatalf("create ogg packer: %s", err)
		}

		opusFilename := "testdata/opus_raw/48k_1ch.opus_raw"
		packets := rawOpusPackets(t, opusFilename)
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

		opusFilename := "testdata/opus_raw/48k_1ch.opus_raw"
		packets := rawOpusPackets(t, opusFilename)
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

func rawOpusPackets(t *testing.T, fname string) [][]byte {
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
