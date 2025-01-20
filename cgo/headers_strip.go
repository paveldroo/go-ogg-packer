package cgo_oggpacker

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
)

func StripHeaders() {
	f, err := os.Open(testFilePath)
	if err != nil {
		log.Fatalf("extract open file: %s", err.Error())
	}

	bufin := bytes.Buffer{}
	bufin.ReadFrom(f)
	b := bufin.Bytes()
	fmt.Println("BUFIN LEN", len(b))

	const headerLen = 27

	bufout := bytes.Buffer{}
	i := 0

	for i < len(b) {
		if reflect.DeepEqual(b[i:i+4], []byte("OggS")) {
			fmt.Println("=========================")
			segmentLen := int(b[i+headerLen])
			fmt.Println("segment len", segmentLen)
			opusDataLength := 0

			for j := 0; j < segmentLen; j++ {
				opusDataLength += int(b[i+headerLen+1+j])
			}

			rawDataStartIdx := i + headerLen + segmentLen + 1
			fmt.Println("raw data start idx", rawDataStartIdx)
			rawDataEndIdx := rawDataStartIdx + opusDataLength
			fmt.Println("raw data end idx", rawDataEndIdx)
			rawOpusData := b[rawDataStartIdx:rawDataEndIdx]
			fmt.Println("ADD RAW OPUS DATA LEN", len(rawOpusData))
			bufout.Write(rawOpusData)

			i = rawDataStartIdx + opusDataLength
		} else {
			i += 1
		}
	}

	fmt.Println("BUFOUT LEN", len(bufout.Bytes()))
	fmt.Println("=========================")

	mustWriteOpusFile("", bufout.Bytes())
}
