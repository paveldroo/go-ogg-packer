package packer_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	packer "github.com/paveldroo/go-ogg-packer"
)

const fileBasePath = "48k_1ch"

func TestPacker(t *testing.T) {
	tests := []struct {
		name        string
		sourceFname string
		refFname    string
		wantErr     bool
		errByte     byte
	}{
		{
			name:        "48k 1ch",
			sourceFname: fmt.Sprintf("testdata/%s.wav", fileBasePath),
			refFname:    fmt.Sprintf("testdata/want/%s.ogg", fileBasePath),
			wantErr:     false,
		},
		{
			name:        "48k 1ch want error",
			sourceFname: fmt.Sprintf("testdata/%s.wav", fileBasePath),
			refFname:    fmt.Sprintf("testdata/want/%s.ogg", fileBasePath),
			wantErr:     true,
			errByte:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcmData := pcmFromWav(t, tt.sourceFname)
			packer, err := packer.New()
			if err != nil {
				t.Fatalf("create new packer: %s", err.Error())
			}

			for i := 0; i < len(pcmData); i++ {
				end := i + 2048
				if end > len(pcmData) {
					end = len(pcmData)
				}
				if err := packer.SendPCMChunk(pcmData[i:end]); err != nil {
					log.Fatalf("send PCM chunk: %s", err.Error())
				}
				i = end
			}

			audioData, err := packer.GetResult()
			if err != nil {
				log.Fatalf("get result from packer: %s", err.Error())
			}

			refData, err := os.ReadFile(tt.refFname)
			if err != nil {
				t.Fatalf("open reference file: %s", err.Error())
			}

			if tt.wantErr {
				audioData = append(audioData, tt.errByte)
				if reflect.DeepEqual(refData, audioData) {
					t.Fatal("source data and want data should NOT be equal")
				}
				return
			}

			if !reflect.DeepEqual(refData, audioData) {
				t.Fatal("source data and want data should NOT be equal")
			}
		})
	}
}

func pcmFromWav(t *testing.T, fn string) []int16 {
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
