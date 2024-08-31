package lib

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/paveldroo/go-ogg-packer/lib/opus_decoder"
)

// OpusFrameDuration is a standard opus frame duration == 20ms == 1 ogg page
const OpusFrameDuration = 20

// SamplesCnt calculates how many bytes fits in one Ogg page
func SamplesCnt(sampleRate int) int {
	return (sampleRate * OpusFrameDuration) / 1000
}

func MustWriteResultFile(data []byte) {
	const resultFilePath = "testdata/audio/result.ogg"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}
	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}

// DecodeOpusOGG not properly working for now
func DecodeOpusOGG() []int16 {
	const opusOGGFilePath = "cgo_oggpacker/testdata/audio/ref/office.opus.ogg"
	decoder := opus_decoder.New(1, 8000)
	d, err := os.ReadFile(opusOGGFilePath)
	if err != nil {
		log.Fatalf("open opus ogg file path: %s", err.Error())
	}

	samplesCount := SamplesCnt(8000)

	var res []int16

	for i := 0; i < len(d); i += samplesCount {
		end := i + samplesCount
		if end > len(d) {
			end = len(d)
		}
		chunk, err := decoder.DecodeInt16(d)
		if err != nil {
			log.Fatalf("decode opus ogg file to Int16: %s", err.Error())
		}

		res = append(res, chunk...)
	}

	return res
}

func MustWriteS16File(data []int16) {
	const resultFilePath = "testdata/audio/office_result.s16"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	buf := make([]byte, len(data)*2)
	for i, v := range data {
		binary.LittleEndian.PutUint16(buf[i*2:], uint16(v))
	}

	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, buf, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}

func ExtractRawOpusFromOGG() []byte {
	const opusOGGFilePath = "cgo_oggpacker/testdata/audio/ref/office.opus.ogg"
	d, err := os.ReadFile(opusOGGFilePath)
	if err != nil {
		log.Fatalf("open opus ogg file path: %s", err.Error())
	}

	pattern := regexp.MustCompile("OggS")
	pageBoundaries := pattern.FindAllIndex(d, -1)
	buf := bytes.Buffer{}
	for i := range pageBoundaries {
		if i < 2 {
			continue
		}
		currPageStartIdx := pageBoundaries[i][0]
		prevOGGDataEnd := pageBoundaries[i-1][0] + 26
		opusData := d[prevOGGDataEnd:currPageStartIdx]
		buf.Write(opusData)
	}
	res := buf.Bytes()
	return res
}

func MustWriteOpusFile(data []byte) {
	const resultFilePath = "cgo_oggpacker/testdata/audio/office_result.opus"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}
