package packer_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/paveldroo/go-ogg-packer/packer"
	"github.com/paveldroo/go-ogg-packer/tests/testutil"
)

const headersCount = 39

// TestPacker compares packer output against reference PCM files.
// References must be generated before running: task generate-ref
// Or run everything together: task test
func TestPacker(t *testing.T) {
	tests := []struct {
		name        string
		sampleRate  int
		channels    int
		sourceFname string
		refFname    string
		wantErr     bool
		errByte     int16
	}{
		{
			name:        "8k 1ch",
			sampleRate:  8000,
			channels:    1,
			sourceFname: "../tests/testdata/input/8k_1ch.pcm",
			refFname:    "../tests/testdata/want/packer/8k_1ch.pcm",
		},
		{
			name:        "12k 1ch",
			sampleRate:  12000,
			channels:    1,
			sourceFname: "../tests/testdata/input/12k_1ch.pcm",
			refFname:    "../tests/testdata/want/packer/12k_1ch.pcm",
		},
		{
			name:        "16k 1ch",
			sampleRate:  16000,
			channels:    1,
			sourceFname: "../tests/testdata/input/16k_1ch.pcm",
			refFname:    "../tests/testdata/want/packer/16k_1ch.pcm",
		},
		{
			name:        "24k 1ch",
			sampleRate:  24000,
			channels:    1,
			sourceFname: "../tests/testdata/input/24k_1ch.pcm",
			refFname:    "../tests/testdata/want/packer/24k_1ch.pcm",
		},
		{
			name:        "48k 1ch",
			sampleRate:  48000,
			channels:    1,
			sourceFname: "../tests/testdata/input/48k_1ch.pcm",
			refFname:    "../tests/testdata/want/packer/48k_1ch.pcm",
		},
		{
			name:        "8k 2ch",
			sampleRate:  8000,
			channels:    2,
			sourceFname: "../tests/testdata/input/8k_2ch.pcm",
			refFname:    "../tests/testdata/want/packer/8k_2ch.pcm",
		},
		{
			name:        "12k 2ch",
			sampleRate:  12000,
			channels:    2,
			sourceFname: "../tests/testdata/input/12k_2ch.pcm",
			refFname:    "../tests/testdata/want/packer/12k_2ch.pcm",
		},
		{
			name:        "16k 2ch",
			sampleRate:  16000,
			channels:    2,
			sourceFname: "../tests/testdata/input/16k_2ch.pcm",
			refFname:    "../tests/testdata/want/packer/16k_2ch.pcm",
		},
		{
			name:        "24k 2ch",
			sampleRate:  24000,
			channels:    2,
			sourceFname: "../tests/testdata/input/24k_2ch.pcm",
			refFname:    "../tests/testdata/want/packer/24k_2ch.pcm",
		},
		{
			name:        "48k 2ch",
			sampleRate:  48000,
			channels:    2,
			sourceFname: "../tests/testdata/input/48k_2ch.pcm",
			refFname:    "../tests/testdata/want/packer/48k_2ch.pcm",
		},
		{
			name:        "48k 1ch want error",
			sampleRate:  48000,
			channels:    1,
			sourceFname: "../tests/testdata/input/48k_1ch.pcm",
			refFname:    "../tests/testdata/want/packer/48k_1ch.pcm",
			wantErr:     true,
			errByte:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcePCMData := testutil.PCMData(t, tt.sourceFname)

			p, err := packer.New(tt.channels, tt.sampleRate, testutil.TestSerialNo)
			if err != nil {
				t.Fatalf("create new packer: %s", err.Error())
			}

			for i := 0; i < len(sourcePCMData); i += packer.DefaultPCMChunkSize {
				end := i + packer.DefaultPCMChunkSize
				if end > len(sourcePCMData) {
					end = len(sourcePCMData)
				}
				if err := p.AddPCMChunk(sourcePCMData[i:end]); err != nil {
					t.Fatalf("send PCM chunk: %s", err.Error())
				}
			}

			audioData, err := p.Result()
			if err != nil {
				log.Fatalf("get result from packer: %s", err.Error())
			}

			pcm := testutil.PCMFromOgg(t, audioData, tt.sampleRate, tt.channels)

			refData := testutil.PCMData(t, tt.refFname)

			if tt.wantErr {
				pcm = append(pcm, tt.errByte)
				if reflect.DeepEqual(refData, pcm) {
					t.Fatal("source data and want data should NOT be equal")
				}
				return
			}

			if mse := testutil.CalculateMSE(t, refData, pcm); mse > 5.0 {
				t.Fatalf("significant distortions in reference and result files, mse: %f", mse)
			}
		})
	}
}
