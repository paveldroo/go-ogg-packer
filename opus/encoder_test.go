package opus_test

import (
	"math/rand"
	"testing"

	"github.com/paveldroo/go-ogg-packer/opus"
)

func TestEncoder_Encode(t *testing.T) {
	cfg := opus.NewDefaultConfig()
	frameSizeSamples := opus.FrameSizeSamples(cfg)

	tests := []struct {
		name       string
		pcmData    []int16
		wantResLen int
		wantPos    int
	}{
		{
			name:       "1x opus packet",
			pcmData:    generateRandomPCMData(frameSizeSamples),
			wantResLen: 1,
			wantPos:    frameSizeSamples,
		},
		{
			name:       "2x opus packet",
			pcmData:    generateRandomPCMData(frameSizeSamples * 2),
			wantResLen: 2,
			wantPos:    frameSizeSamples * 2,
		},
		{
			name:       "0.5x opus packet",
			pcmData:    generateRandomPCMData(int(float32(frameSizeSamples) * 0.5)),
			wantResLen: 0,
			wantPos:    0,
		},
		{
			name:       "2.5x opus packet",
			pcmData:    generateRandomPCMData(int(float32(frameSizeSamples) * 2.5)),
			wantResLen: 2,
			wantPos:    int(float32(frameSizeSamples) * 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder, err := opus.NewEncoder(cfg)
			if err != nil {
				t.Fatalf("create opus encoder: %s", err.Error())
			}

			res, pos, err := encoder.Encode(tt.pcmData)
			if err != nil {
				t.Fatalf("encode pcm data: %s", err.Error())
			}

			if len(res) != tt.wantResLen {
				t.Fatalf("result length should be equal %d, current %d", tt.wantResLen, len(res))
			}

			if pos != tt.wantPos {
				t.Fatalf("position should be equal %d, current %d", tt.wantPos, pos)
			}
		})
	}
}

func TestEncoder_EncodeWithPadding(t *testing.T) {
	cfg := opus.NewDefaultConfig()
	frameSizeSamples := opus.FrameSizeSamples(cfg)

	tests := []struct {
		name       string
		pcmData    []int16
		wantResLen int
	}{
		{
			name:       "1x opus packet",
			pcmData:    generateRandomPCMData(frameSizeSamples),
			wantResLen: 1,
		},
		{
			name:       "1.5x opus packet",
			pcmData:    generateRandomPCMData(int(float32(frameSizeSamples) * 1.5)),
			wantResLen: 2,
		},
		{
			name:       "0.5x opus packet",
			pcmData:    generateRandomPCMData(int(float32(frameSizeSamples) * 0.5)),
			wantResLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder, err := opus.NewEncoder(cfg)
			if err != nil {
				t.Fatalf("create opus encoder: %s", err.Error())
			}

			res, err := encoder.EncodeWithPadding(tt.pcmData)
			if err != nil {
				t.Fatalf("encode pcm data: %s", err.Error())
			}

			if len(res) != tt.wantResLen {
				t.Fatalf("result length should be equal %d, current %d", tt.wantResLen, len(res))
			}
		})
	}
}

func generateRandomPCMData(size int) []int16 {
	pcm := make([]int16, size)
	for i := range pcm {
		pcm[i] = int16(rand.Intn(65536) - 32768)
	}
	return pcm
}
