package tests

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	packer "github.com/paveldroo/go-ogg-packer"
	writer "github.com/paveldroo/go-ogg-packer/tests/buffer_writer"
	"github.com/paveldroo/go-ogg-packer/tests/buffer_writer/opus_tools"
)

const baseFilename = "testdata/base.ogg"

func TestPacker(t *testing.T) {
	converter, err := opus_tools.NewOpusConverter(opus_tools.NewDefaultConfig())
	if err != nil {
		t.Fatalf("create opus converter: %s", err.Error())
	}

	packer, err := packer.New(1, sampleRate)
	if err != nil {
		t.Fatalf("create ogg packer: %s", err.Error())
	}

	s16 := s16FromWav()
	audioBuffer := writer.NewAudioBuffer(converter, packer)

	for i := 0; i < len(s16); i++ {
		end := i + 2048
		if end > len(s16) {
			end = len(s16)
		}
		if err := audioBuffer.SendS16Chunk(s16[i:end]); err != nil {
			t.Fatalf("send s16 chunk: %s", err.Error())
		}
		i = end
	}

	audioContent, err := audioBuffer.GetResult()
	if err != nil {
		t.Fatalf("get result from audio buffer: %s", err.Error())
	}

	fname := fmt.Sprintf("testdata/result/ogg_packer_result_%d.ogg", time.Now().UnixNano())
	mustWriteOggFile(fname, audioContent)

	baseData, err := os.ReadFile(baseFilename)
	if err != nil {
		t.Fatalf("open base file: %s", err.Error())
	}

	if !reflect.DeepEqual(baseData, audioContent) {
		t.Fatal("base data and test data are not equal")
	}
}
