package opus

import "errors"

// frameDuration48kHz maps opus TOC config (0–31) to frame duration
// in 48kHz samples per channel (RFC 6716 §3.1).
var frameDuration48kHz = [32]int{
	// SILK-only (NB/MB/WB × 10/20/40/60 ms)
	480, 960, 1920, 2880,
	480, 960, 1920, 2880,
	480, 960, 1920, 2880,
	// Hybrid (SWB/FB × 10/20 ms)
	480, 960,
	480, 960,
	// CELT-only (NB/WB/SWB/FB × 2.5/5/10/20 ms)
	120, 240, 480, 960,
	120, 240, 480, 960,
	120, 240, 480, 960,
	120, 240, 480, 960,
}

// SamplesCount parses the opus TOC byte (RFC 6716 §3.1) and returns
// the total number of PCM samples (across all channels) at the given sample rate.
func SamplesCount(data []byte, channelCount uint8, sampleRate uint32) (int, error) {
	if len(data) < 1 {
		return 0, errors.New("empty opus packet")
	}

	toc := data[0]
	config := toc >> 3
	code := toc & 0x03

	frameDuration := frameDuration48kHz[config]

	var frameCount int
	switch code {
	case 0:
		frameCount = 1
	case 1, 2:
		frameCount = 2
	case 3:
		if len(data) < 2 {
			return 0, errors.New("opus code 3 packet too short")
		}
		frameCount = int(data[1] & 0x3F)
	}

	samplesPerChannel := frameDuration * frameCount * int(sampleRate) / 48000
	return samplesPerChannel * int(channelCount), nil
}
