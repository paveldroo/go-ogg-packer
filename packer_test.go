package packer_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"testing"

	extopus "gopkg.in/hraban/opus.v2"
	extogg "mccoy.space/g/ogg"

	"github.com/google/go-cmp/cmp"
	packer "github.com/paveldroo/go-ogg-packer"
	"github.com/paveldroo/go-ogg-packer/opus"
)

const (
	fileBasePath = "48k_1ch"
	headersCount = 39
)

func TestPacker(t *testing.T) {
	tests := []struct {
		name        string
		sourceFname string
		refFname    string
		wantErr     bool
		errByte     int16
	}{
		{
			name:        "48k 1ch",
			sourceFname: fmt.Sprintf("testdata/%s.pcm", fileBasePath),
			refFname:    fmt.Sprintf("testdata/want/%s.pcm", fileBasePath),
			wantErr:     false,
		},
		{
			name:        "48k 1ch want error",
			sourceFname: fmt.Sprintf("testdata/%s.pcm", fileBasePath),
			refFname:    fmt.Sprintf("testdata/want/%s.pcm", fileBasePath),
			wantErr:     true,
			errByte:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcePCMData := pcmData(t, tt.sourceFname)
			refData := pcmData(t, tt.refFname)

			packer, err := packer.New()
			if err != nil {
				t.Fatalf("create new packer: %s", err.Error())
			}

			for i := 0; i < len(sourcePCMData); i++ {
				end := min(i+2048, len(sourcePCMData))
				if err := packer.SendPCMChunk(sourcePCMData[i:end]); err != nil {
					t.Fatalf("send PCM chunk: %s", err.Error())
				}
				i = end
			}

			audioData, err := packer.GetResult()
			if err != nil {
				log.Fatalf("get result from packer: %s", err.Error())
			}

			pcm := pcmFromOgg(t, audioData)

			if tt.wantErr {
				pcm = append(pcm, tt.errByte)
				if diff := cmp.Diff(refData, pcm, TolerantByteDiff(2)); diff == "" {
					t.Fatal("source data and want data should NOT be equal")
				}
				return
			}

			if diff := cmp.Diff(refData, pcm, TolerantByteDiff(2)); diff != "" {
				t.Fatal("source data and want data should be equal with acceptable tolerance", diff)
			}
		})
	}
}

func pcmData(t *testing.T, fn string) []int16 {
	t.Helper()

	d, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("open wav file: %s", err.Error())
	}

	reader := bytes.NewReader(d)
	numValues := len(d) / 2

	result := make([]int16, numValues)

	for i := range result {
		var value int16
		if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
			t.Fatalf("binary read wav file: %s", err.Error())
		}
		result[i] = value
	}

	return result
}

func pcmFromOgg(t *testing.T, oggData []byte) []int16 {
	t.Helper()

	b := bytes.NewBuffer(oggData)
	oggDecoder := extogg.NewDecoder(b)

	opusDecoder, err := extopus.NewDecoder(opus.SampleRate, opus.NumChannels)
	if err != nil {
		t.Fatalf("create opus decoder: %s", err.Error())
	}

	pcmBuffer := make([]int16, opus.FrameSize*opus.SampleRate*opus.NumChannels/1000)

	var pcm []int16
	for {
		page, err := oggDecoder.Decode()
		if err != nil {
			break
		}

		for _, packet := range page.Packets {
			n, err := opusDecoder.Decode(packet, pcmBuffer)
			if err != nil {
				continue // some errors are acceptable during packet decoding
			}
			pcm = append(pcm, pcmBuffer[:n]...)
		}
	}

	return pcm
}

// TolerantByteDiff returns a cmp.Option that allows a small tolerance when comparing []byte.
func TolerantByteDiff(tolerance int) cmp.Option {
	return cmp.FilterValues(func(x, y []byte) bool {
		return len(x) == len(y)
	}, cmp.Comparer(func(x, y []byte) bool {
		for i := 0; i < len(x); i++ {
			if int(math.Abs(float64(int(x[i])-int(y[i])))) > tolerance {
				return false
			}
		}
		return true
	}))
}
