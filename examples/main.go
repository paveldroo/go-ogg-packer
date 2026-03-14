package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	packer "github.com/paveldroo/go-ogg-packer"
	"github.com/paveldroo/go-ogg-packer/opus"
)

const (
	wavFilePath = "examples/12k_1ch.wav"
	sampleRate  = 12000
)

func main() {
	pcmData := pcmFromWav()

	cfg := opus.Config{
		SampleRate:  sampleRate,
		NumChannels: opus.NumChannels,
		FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
	}
	packer, err := packer.New(cfg)
	if err != nil {
		log.Fatalf("create new packer: %s", err.Error())
	}

	const chunkSize = 2048
	for i := 0; i < len(pcmData); i += chunkSize {
		end := min(i+chunkSize, len(pcmData))
		if err := packer.SendPCMChunk(pcmData[i:end]); err != nil {
			log.Fatalf("send s16 chunk: %s", err.Error())
		}
	}

	audioContent, err := packer.GetResult()
	if err != nil {
		log.Fatalf("get result from packer: %s", err.Error())
	}

	fname := fmt.Sprintf("examples/packer_result_%d.ogg", time.Now().UnixNano())
	if err := writeOggFile(fname, audioContent); err != nil {
		log.Fatalf("write ogg file: %s", err.Error())
	}
}

func pcmFromWav() []int16 {
	d, err := os.ReadFile(wavFilePath)
	if err != nil {
		log.Fatalf("open wav file: %s", err.Error())
	}

	wavData := stripWithOffset(d)
	reader := bytes.NewReader(wavData)
	numValues := len(wavData) / 2

	result := make([]int16, numValues)

	for i := range result {
		var value int16
		if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
			log.Fatalf("binary read wav file: %s", err.Error())
		}
		result[i] = value
	}

	return result
}

func stripWithOffset(wavData []byte) []byte {
	// Find the "data" chunk — WAV headers can vary in size due to extra chunks.
	dataOffset := -1
	pos := 12 // skip RIFF header
	for pos < len(wavData)-8 {
		chunkID := string(wavData[pos : pos+4])
		chunkSize := int(binary.LittleEndian.Uint32(wavData[pos+4 : pos+8]))
		if chunkID == "data" {
			dataOffset = pos + 8
			break
		}
		pos += 8 + chunkSize
		if chunkSize%2 != 0 {
			pos++ // align to even boundary
		}
	}
	if dataOffset < 0 {
		log.Fatalf("no data chunk found in WAV file")
	}
	return wavData[dataOffset:]
}

func writeOggFile(name string, data []byte) error {
	wDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current work directory: %w", err)
	}

	var fPath = path.Join(wDir, name)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		return fmt.Errorf("write result file: %w", err)
	}

	return nil
}
