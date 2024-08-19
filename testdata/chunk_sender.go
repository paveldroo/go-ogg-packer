package testdata

import (
	"log"
	"os"
	"path"
)

// NewChunkSender sends file in channel chunk by chunk
func NewChunkSender(fPath string, chunkSize int) <-chan []byte {
	d, err := os.ReadFile(fPath)
	if err != nil {
		log.Fatalf("open file in NewChunkSender: %s", err)
	}

	ch := make(chan []byte)
	go func() {
		for i := 0; i < len(d); i += chunkSize {
			if i+chunkSize > len(d) {
				chunkSize = len(d) - i
			}
			ch <- d[i:chunkSize]
		}
	}()

	return ch
}

func MustReferencePath() string {
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get working directory: %s", err)
	}

	return path.Join(wDir, "audio/ref/office.ogg")
}
