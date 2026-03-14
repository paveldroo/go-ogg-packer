package opus

import "testing"

func TestSamplesCount(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		channels   uint8
		sampleRate uint32
		want       int
		wantErr    bool
	}{
		// SILK config=0 (10ms, 480 samples/ch @48k), code=0 (1 frame), mono
		{name: "SILK 10ms 1fr mono", data: []byte{0b00000_0_00}, channels: 1, sampleRate: 48000, want: 480},
		// SILK config=2 (40ms, 1920 samples/ch @48k), code=0 (1 frame), stereo
		{name: "SILK 40ms 1fr stereo", data: []byte{0b00010_0_00}, channels: 2, sampleRate: 48000, want: 3840},
		// SILK config=3 (60ms, 2880 samples/ch @48k), code=1 (2 frames), mono
		{name: "SILK 60ms 2fr mono", data: []byte{0b00011_0_01}, channels: 1, sampleRate: 48000, want: 5760},
		// Hybrid config=12 (10ms, 480 samples/ch @48k), code=0 (1 frame), mono
		{name: "Hybrid 10ms 1fr mono", data: []byte{0b01100_0_00}, channels: 1, sampleRate: 48000, want: 480},
		// Hybrid config=13 (20ms, 960 samples/ch @48k), code=2 (2 frames), stereo
		{name: "Hybrid 20ms 2fr stereo", data: []byte{0b01101_0_10}, channels: 2, sampleRate: 48000, want: 3840},
		// CELT config=16 (2.5ms, 120 samples/ch @48k), code=0 (1 frame), mono
		{name: "CELT 2.5ms 1fr mono", data: []byte{0b10000_0_00}, channels: 1, sampleRate: 48000, want: 120},
		// CELT config=19 (20ms, 960 samples/ch @48k), code=0 (1 frame), mono
		{name: "CELT 20ms 1fr mono", data: []byte{0b10011_0_00}, channels: 1, sampleRate: 48000, want: 960},
		// CELT config=19 (20ms, 960 samples/ch @48k), code=1 (2 frames), stereo
		{name: "CELT 20ms 2fr stereo", data: []byte{0b10011_0_01}, channels: 2, sampleRate: 48000, want: 3840},
		// Code 3: frame count from data[1] & 0x3F
		{name: "code3 3frames mono", data: []byte{0b00000_0_11, 0x03}, channels: 1, sampleRate: 48000, want: 1440},
		{name: "code3 5frames stereo", data: []byte{0b10011_0_11, 0x05}, channels: 2, sampleRate: 48000, want: 9600},
		// Non-48kHz: SILK config=3 (60ms), code=0 (1 frame), mono @16kHz
		// 2880 * 16000 / 48000 = 960
		{name: "SILK 60ms 1fr mono 16k", data: []byte{0b00011_0_00}, channels: 1, sampleRate: 16000, want: 960},
		// Edge cases
		{name: "empty data", data: []byte{}, channels: 1, sampleRate: 48000, wantErr: true},
		{name: "code3 too short", data: []byte{0b00000_0_11}, channels: 1, sampleRate: 48000, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SamplesCount(tt.data, tt.channels, tt.sampleRate)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
