package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	packer "github.com/paveldroo/go-ogg-packer"
)

const (
	channelCount    = 1
	sampleRate      = 48000
	rawOpusFilename = "48k_1ch.opus_raw"
)

func main() {
	packer, err := packer.New(uint8(channelCount), uint32(sampleRate))
	if err != nil {
		log.Fatalf("create ogg packer: %s", err.Error())
	}

	rawOpusData, err := getRawOpusPackets()
	if err != nil {
		log.Fatalf("get raw opus data: %s", err.Error())
	}

	for _, packet := range rawOpusData {
		if err := packer.AddChunk(packet, false, -1); err != nil {
			log.Fatalf("send opus chunk to packer: %s", err.Error())
		}
	}

	oggData, err := packer.ReadPages()
	if err != nil {
		log.Fatalf("read all pages from packer: %s", err.Error())
	}

	fname := fmt.Sprintf("ogg_packer_result_%d.ogg", time.Now().UnixNano())
	writeOggFile(fname, oggData)
}

func getRawOpusPackets() ([][]byte, error) {
	f, err := os.Open(rawOpusFilename)
	if err != nil {
		return nil, fmt.Errorf("read raw opus file: %w", err)
	}
	decoder := gob.NewDecoder(f)
	var audioData [][]byte
	if err := decoder.Decode(&audioData); err != nil {
		return nil, fmt.Errorf("decode data from file: %w", err)
	}

	return audioData, nil
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
