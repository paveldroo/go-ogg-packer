package cgo_oggpacker

import (
	"fmt"
	"testing"
)

func Test_OGGPacker(t *testing.T) {
	wrapper, err := newOggPackerWrapper(sampleRate, 1)
	if err != nil {
		t.Fatalf("create ogg packer wrapper: %s", err.Error())
	}

	// const opusOGGFilePath = "../lib/cgo_oggpacker/testdata/audio/ref/demo_48k_1ch.opus"
	// d, err := os.ReadFile(opusOGGFilePath)
	// if err != nil {
	// 	t.Fatalf("open opus ogg file path: %s", err.Error())
	// }

	// chunk := []byte{}

	// for _, b := range d {
	// 	if len(chunk) == 20 {
	// 		wrapper.addChunk(chunk)
	// 		chunk = []byte{}
	// 	}

	// 	chunk = append(chunk, b)
	// }
	//


	НАДО ПРИДУМАТЬ КАК ПОЛУЧАТЬ ЧИСТЫЙ СТРИМ OPUS, СЕЙЧАС СЫПЕМСЯ ИЗЗА ТОГО ЧТО ПЫТАЕМСЯ ЗАСУНУТЬ OGG ФАЙЛ


	d := AudioByChunks()

	for _, chunk := range d {
		wrapper.addChunk(chunk)
	}

	res, err := wrapper.readAudioData()
	if err != nil {
		t.Fatalf("read audio data: %s", err.Error())
	}

	fmt.Println("res", res)

}
