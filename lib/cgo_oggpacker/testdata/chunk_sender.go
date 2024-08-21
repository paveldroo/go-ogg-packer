package testdata

import (
	"log"
	"os"
	"path"

	"github.com/paveldroo/go-ogg-packer/lib"
)

const sampleRate = 8000

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
	const testFilePath = "testdata/audio/ref/office.opus.ogg"
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
