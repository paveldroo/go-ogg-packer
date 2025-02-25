package cgo_oggpacker

import (
	"fmt"
	"testing"
	"time"

	"github.com/paveldroo/go-ogg-packer/opus"
	"gitlab.tcsbank.ru/speech/libanysound/ogg_packer"
)

func TestOGGPacker(t *testing.T) {
	converter, err := opus.NewOpusConverter(opus.NewDefaultConfig())
	if err != nil {
		t.Fatalf("create opus converter: %s", err.Error())
	}

	packer, err := ogg_packer.New(sampleRate, 1)
	if err != nil {
		t.Fatalf("create ogg packer wrapper: %s", err.Error())
	}

	s16 := S16FromWav()
	pcm, err := Int16SliceToByteSlice(s16)
	if err != nil {
		t.Fatalf("convert int16 to byte: %s", err.Error())
	}

	mustWriteS16File(pcm)

	audioBuffer := NewAudioBuffer(converter, packer)

	for i := 0; i < len(s16); i++ {
		end := i + 2048
		if end > len(s16) {
			end = len(s16)
		}
		audioBuffer.sendS16Chunk(s16[i:end])
		i = end
	}

	audioContent, err := audioBuffer.getResult()
	if err != nil {
		t.Fatalf("get result from audio buffer: %s", err.Error())
	}

	mustWriteOpusFile(fmt.Sprintf("testdata/result/ogg_packer_result_%d.opus", time.Now().UnixNano()), audioContent)
}
