package testdata

import (
	"encoding/binary"
	"fmt"
	"io"
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
	const testFilePath = "testdata/audio/ref/demo_48k_1ch.opus"
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

func rawPCM() []byte {
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current work directory: %s", err.Error())
	}
	file, err := os.Open(path.Join(wDir, "testdata/audio/raw_pcm_48k_1ch.dat"))
	if err != nil {
		log.Fatalf("open file: %s", err.Error())
	}
	defer file.Close()

	var samples []byte

	for {
		var sample int16
		err := binary.Read(file, binary.LittleEndian, &sample)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read one sample from file: %s", err.Error())
		}
		sampleBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(sampleBytes, uint16(sample))
		samples = append(samples, sampleBytes...)
	}

	fmt.Println("Number of samples read:", len(samples))
	fmt.Println("First 10 samples:", samples[:20])

	return samples
}
