package cgo_oggpacker

import (
	"bytes"
	"fmt"

	"gitlab.tcsbank.ru/speech/libanysound/ogg_packer"
)

type oggPackerWrapper struct {
	packer      *ogg_packer.Packer
	audioBuffer bytes.Buffer
}

func newOggPackerWrapper(sampleRate int, numChannels int) (*oggPackerWrapper, error) {
	oggPacker, err := ogg_packer.New(sampleRate, numChannels)
	if err != nil {
		return nil, fmt.Errorf("create ogg packer: %w", err)
	}
	return &oggPackerWrapper{
		packer:      oggPacker,
		audioBuffer: bytes.Buffer{},
	}, nil
}

func (packerWrapper *oggPackerWrapper) addChunk(chunk []byte) error {
	if packerWrapper.audioBuffer.Len() > 0 {
		previousChunk := packerWrapper.audioBuffer.Bytes()
		err := packerWrapper.packer.AddChunk(previousChunk, false, 960)
		if err != nil {
			return fmt.Errorf("add chunk to ogg packer: %w", err)
		}
		packerWrapper.audioBuffer.Reset()
	}
	_, err := packerWrapper.audioBuffer.Write(chunk)
	if err != nil {
		return fmt.Errorf("add chunk to buffer: %w", err)
	}
	return nil
}

func (packerWrapper *oggPackerWrapper) readAudioData() ([]byte, error) {
	if packerWrapper.audioBuffer.Len() > 0 {
		lastChunk := packerWrapper.audioBuffer.Bytes()
		fmt.Println("lastChunk len", len(lastChunk))
		err := packerWrapper.packer.AddChunk(lastChunk, true, len(lastChunk))
		if err != nil {
			return nil, fmt.Errorf("add chunk to ogg packer: %w", err)
		}
	}

	oggPages, err := packerWrapper.packer.ReadPages()
	if err != nil {
		return nil, fmt.Errorf("get pages from ogg packer: %w", err)
	}
	flushedPages, err := packerWrapper.packer.FlushPages()
	if err != nil {
		return nil, fmt.Errorf("flush pages from ogg packer: %w", err)
	}
	oggPages = append(oggPages, flushedPages...)
	return oggPages, nil
}

func (packerWrapper *oggPackerWrapper) isEmpty() bool {
	return packerWrapper.audioBuffer.Len() <= 0
}

func (packerWrapper *oggPackerWrapper) close() {
	packerWrapper.packer.Close()
}
