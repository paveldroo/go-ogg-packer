package cgo_oggpacker

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path"
)

const sampleRate = 48000
const channels = 1
const wavFilePath = "testdata/demo_48k_1ch.wav"

func S16FromWav() []int16 {
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

func mustWriteOpusFile(name string, data []byte) {
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, name)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}

func mustWriteWavFile(data []byte) {
	const resultFilePath = "testdata/demo_result.wav"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}

func mustWriteS16File(data []byte) {
	const resultFilePath = "testdata/result/int16.pcm"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}

func Int16SliceToByteSlice(int16s []int16) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, val := range int16s {
		err := binary.Write(buf, binary.LittleEndian, val)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
