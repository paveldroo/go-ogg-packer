package cgo_oggpacker

import (
	"log"
	"os"
	"path"

	"github.com/paveldroo/go-ogg-packer/lib"
)

const sampleRate = 48000

func AudioByChunks() [][]byte {
	var chunkSize = lib.SamplesCnt(sampleRate)
	d := RefOGGData()

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

func RefOGGData() []byte {
	const testFilePath = "../lib/cgo_oggpacker/testdata/audio/ref/demo_48k_1ch.opus"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}
	var fPath = path.Join(wDir, testFilePath)

	d, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatalf("open file in RefOGGData: %s", err)
	}

	return d
}

func mustWriteOpusFile(data []byte) {
	const resultFilePath = ".cgo_oggpacker/testdata/audio/office_result.opus"
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}

	var fPath = path.Join(wDir, resultFilePath)
	if err := os.WriteFile(fPath, data, 0666); err != nil {
		log.Fatalf("write result file: %s", err.Error())
	}
}
