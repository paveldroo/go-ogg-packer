package lib_test

import (
	"testing"

	"github.com/paveldroo/go-ogg-packer/lib"
)

func TestExtractRawOpusFromOGG(t *testing.T) {
	opusData := lib.ExtractRawOpusFromOGG()
	lib.MustWriteOpusFile(opusData)
}
