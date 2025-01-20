package cgo_oggpacker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/hraban/opus"
)

const sampleRate = 48000
const channels = 1
const testFilePath = "testdata/demo_48k_1ch.opus"

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
	// d := RefOGGData()
	d := ExtractOpusFromOGG()

	var res [][]byte

	for i := 0; i < len(d); i += chunkSize {
		end := i + chunkSize
		if end > len(d) {
			end = len(d)
		}
		res = append(res, d[i:end])
	}

	fmt.Println("bytes result", res)
	fmt.Println("bytes len", len(res))

	return res
}

func mustWriteOpusFile(data []byte) {
	const resultFilePath = "testdata/office_result.opus"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}

func mustWriteWavFile(data []byte) {
	const resultFilePath = "testdata/office_result.wav"
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
