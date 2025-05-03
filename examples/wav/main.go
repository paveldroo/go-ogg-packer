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
	writer "github.com/paveldroo/go-ogg-packer/examples/wav/buffer_writer"
	"github.com/paveldroo/go-ogg-packer/examples/wav/buffer_writer/opus_tools"
)

const (
	sampleRate  = 48000
	wavFilePath = "48k_1ch.wav"
)

func main() {
	converter, err := opus_tools.NewOpusConverter(opus_tools.NewDefaultConfig())
	if err != nil {
		log.Fatalf("create opus converter: %s", err.Error())
	}

	packer, err := packer.New(1, sampleRate)
	if err != nil {
		log.Fatalf("create ogg packer: %s", err.Error())
	}

	s16 := s16FromWav()
	audioBuffer := writer.NewAudioBuffer(converter, packer)

	for i := 0; i < len(s16); i++ {
		end := i + 2048
		if end > len(s16) {
			end = len(s16)
		}
		if err := audioBuffer.SendS16Chunk(s16[i:end]); err != nil {
			log.Fatalf("send s16 chunk: %s", err.Error())
		}
		i = end
	}

	audioContent, err := audioBuffer.GetResult()
	if err != nil {
		log.Fatalf("get result from audio buffer: %s", err.Error())
	}

	fname := fmt.Sprintf("ogg_packer_result_%d.ogg", time.Now().UnixNano())
	if err := writeOggFile(fname, audioContent); err != nil {
		log.Fatalf("write ogg file: %s", err.Error())
	}
}

func s16FromWav() []int16 {
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
