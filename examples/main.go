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
)

const wavFilePath = "examples/48k_1ch.wav"

func main() {
	pcmData := pcmFromWav()
	packer, err := packer.New()
	if err != nil {
		log.Fatalf("create new packer: %s", err.Error())
	}

	for i := 0; i < len(pcmData); i++ {
		end := i + 2048
		if end > len(pcmData) {
			end = len(pcmData)
		}
		if err := packer.SendPCMChunk(pcmData[i:end]); err != nil {
			log.Fatalf("send s16 chunk: %s", err.Error())
		}
		i = end
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

	reader := bytes.NewReader(d)
	numValues := len(d) / 2

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
