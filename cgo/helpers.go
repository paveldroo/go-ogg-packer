package cgo_oggpacker

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/hraban/opus.v2"
)

const sampleRate = 48000
const channels = 1
const testFilePath = "testdata/raw_opus.opus"
const wavFilePath = "testdata/demo_48k_1ch_raw_opus.wav"

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

func ExtractOpusFromOGG() []byte {
	f, err := os.Open(testFilePath)
	if err != nil {
		log.Fatalf("extract open file: %s", err.Error())
	}
	s, err := opus.NewStream(f)
	if err != nil {
		log.Fatalf("create opus stream: %s", err.Error())
	}
	defer s.Close()
	buf := make([]int16, 16384)
	res := []int16{}
	for {
		n, err := s.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("write int16 to buf: %s", err.Error())
		}
		pcm := buf[:n*channels]
		res = append(res, pcm...)
	}

	b, err := Int16SliceToByteSlice(res)
	if err != nil {
		log.Fatalf("int16 to bytes: %s", err.Error())
	}

	mustWriteWavFile(b)

	return b
}

func AudioByChunks() [][]byte {
	var chunkSize = SamplesCnt(sampleRate)
	d, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("audio by chunks: open file: %s", err.Error())
	}

	var res [][]byte

	for i := 0; i < len(d); i += chunkSize {
		end := i + chunkSize
		if end > len(d) {
			end = len(d)
		}
		res = append(res, d[i:end])
	}

	return res
}

func AudioFull() []byte {
	d, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("audio full: open file: %s", err.Error())
	}

	return d
}

func mustWriteOpusFile(name string, data []byte) {
	if name == "" {
		name = "testdata/raw_opus.opus"
	}
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

// OpusFrameDuration is a standard opus frame duration == 20ms == 1 ogg page
const OpusFrameDuration = 20

// SamplesCnt calculates how many bytes fits in one Ogg page
func SamplesCnt(sampleRate int) int {
	return (sampleRate * OpusFrameDuration) / 1000
}
