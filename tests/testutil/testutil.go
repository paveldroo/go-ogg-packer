package testutil

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"testing"

	extopus "gopkg.in/hraban/opus.v2"
	extogg "mccoy.space/g/ogg"

	"github.com/paveldroo/go-ogg-packer/opus"
)

// PCMData reads a raw PCM file and returns int16 samples.
func PCMData(t testing.TB, fn string) []int16 {
	t.Helper()

	d, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("open pcm file: %s", err)
	}

	reader := bytes.NewReader(d)
	result := make([]int16, len(d)/2)

	for i := range result {
		if err := binary.Read(reader, binary.LittleEndian, &result[i]); err != nil {
			t.Fatalf("binary read pcm file: %s", err)
		}
	}

	return result
}

// RawOpusPackets reads a GOB-encoded file of opus packets.
func RawOpusPackets(t testing.TB, fname string) [][]byte {
	t.Helper()

	f, err := os.Open(fname)
	if err != nil {
		t.Fatalf("read raw opus file: %s", err)
	}
	defer f.Close()

	var audioData [][]byte
	if err := gob.NewDecoder(f).Decode(&audioData); err != nil {
		t.Fatalf("decode data from file: %s", err)
	}

	return audioData
}

// WriteOggFile writes binary data to an OGG file.
func WriteOggFile(t testing.TB, fname string, data []byte) {
	t.Helper()

	wDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current work directory: %s", err)
	}

	fPath := path.Join(wDir, fname)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		t.Fatalf("write result file: %s", err)
	}
}

// GenNewRef writes int16 PCM samples to a reference file.
func GenNewRef(t testing.TB, refFileName string, pcmData []int16) {
	t.Helper()

	file, err := os.Create(refFileName)
	if err != nil {
		t.Fatalf("create reference file: %s", err)
	}
	defer file.Close()

	for _, sample := range pcmData {
		if err := binary.Write(file, binary.LittleEndian, sample); err != nil {
			t.Fatalf("write pcm data to file: %s", err)
		}
	}

	fmt.Printf("New reference file %s successfully generated\n", refFileName)
}

// PCMFromOgg decodes OGG/Opus data back to PCM samples.
func PCMFromOgg(t testing.TB, oggData []byte, sampleRate, numChannels int) []int16 {
	t.Helper()

	b := bytes.NewBuffer(oggData)
	oggDecoder := extogg.NewDecoder(b)

	opusDecoder, err := extopus.NewDecoder(sampleRate, numChannels)
	if err != nil {
		t.Fatalf("create opus decoder: %s", err)
	}

	pcmBuffer := make([]int16, opus.FrameSize*sampleRate*numChannels/1000)

	var pcm []int16
	for {
		page, err := oggDecoder.Decode()
		if err != nil {
			break
		}

		for _, packet := range page.Packets {
			n, err := opusDecoder.Decode(packet, pcmBuffer)
			if err != nil {
				continue
			}
			pcm = append(pcm, pcmBuffer[:n]...)
		}
	}

	return pcm
}

// CalculateMSE returns Mean Squared Error between two PCM signals.
func CalculateMSE(t testing.TB, ref, pcm []int16) float64 {
	t.Helper()

	if len(ref) != len(pcm) {
		t.Fatalf("reference and result files lengths not equal")
	}

	var sumSq float64
	for i := range ref {
		diff := int64(ref[i]) - int64(pcm[i])
		sumSq += float64(diff * diff)
	}

	return sumSq / float64(len(ref))
}
