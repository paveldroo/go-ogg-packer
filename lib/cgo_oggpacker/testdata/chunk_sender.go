package testdata

import (
	"log"
	"os"
	"path"
)

func AudioByChunks() [][]byte {
	var chunkSize = 512
	d := RefOGGData()

	var res [][]byte

	for i := 0; i < len(d); i += chunkSize {
		end := i + chunkSize
		if end > len(d) {
			end = len(d) - 1
		}
		res = append(res, d[i:end])
	}

	return res
}

// NewRefAudioChunkSender sends reference audio file in channel chunk by chunk
//func NewRefAudioChunkSender() <-chan []byte {
//	const chunkSize = 512
//	d := RefOGGData()
//	ch := make(chan []byte)
//	go func() {
//		for i := 0; i < len(d); i += chunkSize {
//			end := i+chunkSize
//			if end > len(d) {
//				end = len(d) - 1
//			}
//			ch <- d[i:end]
//		}
//	}()
//
//	return ch
//}

func RefOGGData() []byte {
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}
	var fPath = path.Join(wDir, "testdata/audio/ref/office.ogg")

	d, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatalf("open file in RefOGGData: %s", err)
	}

	return d
}
