package packer_test

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	oggPacker "github.com/paveldroo/go-ogg-packer"
)

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
			sourceFname: "testdata/48k_1ch.pcm",
			refFname:    "testdata/want/48k_1ch.ogg",
			wantErr:     false,
		},
		{
			name:        "48k 1ch want error",
			sourceFname: "testdata/48k_1ch.pcm",
			refFname:    "testdata/want/48k_1ch.ogg",
			wantErr:     true,
			errByte:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcmData := pcmFromFile(t, tt.sourceFname)
			packer, err := oggPacker.New()
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

func TestCustomPacker(t *testing.T) {
	tests := []struct {
		name        string
		sourceFname string
		refFname    string
		wantErr     bool
		errByte     byte
	}{
		{
			name:        "24k 1ch",
			sourceFname: "testdata/24k_1ch.pcm",
			refFname:    "testdata/want/24k_1ch.ogg",
			wantErr:     false,
		},
		{
			name:        "24k 1ch want error",
			sourceFname: "testdata/24k_1ch.pcm",
			refFname:    "testdata/want/24k_1ch.ogg",
			wantErr:     true,
			errByte:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcmData := pcmFromFile(t, tt.sourceFname)
			packer, err := oggPacker.NewWithConfig(24000, 1, 60*time.Millisecond)
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
				t.Fatal("source data and want data should be equal")
			}
		})
	}
}

func TestResample(t *testing.T) {
	rawData := pcmFromFile(t, "testdata/24k_1ch.pcm")
	packer, err := oggPacker.New()
	if err != nil {
		t.Fatalf("create new packer: %s", err.Error())
	}
	pcmData := oggPacker.ResampleLinearInt16(rawData, 24000, 48000)
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
	os.WriteFile("test/test_resample.ogg", audioData, 0666)
}

func pcmFromFile(t *testing.T, fn string) []int16 {
	t.Helper()

	d, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("open audio file: %s", err.Error())
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
