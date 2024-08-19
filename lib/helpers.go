package lib

import (
	"log"
	"os"
	"path"
)

// OpusFrameDuration is a standard opus frame duration == 20ms == 1 ogg page
const OpusFrameDuration = 20

// SamplesCnt calculates how many bytes fits in one Ogg page
func SamplesCnt(sampleRate int) int {
	return (sampleRate * OpusFrameDuration) / 1000
}

func MustWriteResultFile(data []byte) {
	const resultFilePath = "testdata/audio/office_result.ogg"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}
	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}
