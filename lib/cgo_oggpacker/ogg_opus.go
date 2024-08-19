// Package cgo_oggpacker is a CGO ogg packer implementation which a new native Go implementation will be tested against.
package cgo_oggpacker

import "C"
import (
	"errors"
	"math/rand"
	"time"
	"unsafe"
)

/*
#cgo pkg-config: opus ogg
#include "ogg_opus_packer.h"
*/
import "C"

type Packer struct {
	c_object *C.ogg_opus_packer_t
}

func New(sampleRate, numChannels int) (*Packer, error) {
	oggPacker := C.ogg_opus_packer_create()
	if oggPacker == nil {
		return nil, errors.New("Failed to create OGG packer")
	}
	serialNo := rand.New(rand.NewSource(time.Now().UTC().Unix() % 0x80000000)).Int31()
	status := C.ogg_opus_packer_init(oggPacker, C.uint8_t(numChannels), C.uint32_t(sampleRate), C.int(serialNo))
	if err := initStatusToError(int(status)); err != nil {
		return nil, err
	}
	return &Packer{c_object: oggPacker}, nil
}

/*If number of samples is unknow samplesCount < 0*/
func (packer *Packer) AddChunk(data []byte, eos bool, samplesCount int) error {
	eosNumber := 0
	if eos {
		eosNumber = 1
	}
	status := C.ogg_opus_packer_add_opus_chunk(packer.c_object,
		unsafe.Pointer(&data[0]), C.size_t(len(data)), C.int(eosNumber), C.int(samplesCount))
	if status == -1 {
		return errors.New("Failed to add chunk to the stream")
	} else if status < -1 {
		return errors.New("Failed to decode opus chunck:" + opusDecoderStatusToError(int(status)+1).Error())
	}
	return nil
}

func (packer *Packer) ReadPages() ([]byte, error) {
	if status := C.ogg_opus_packer_collect_pages(packer.c_object); status == -1 {
		return nil, errors.New("Failed to add chunk to buffer")
	}
	return packer.readBuffer(), nil
}

func (packer *Packer) FlushPages() ([]byte, error) {
	if status := C.ogg_opus_packer_flush_pages(packer.c_object); status == -1 {
		return nil, errors.New("Failed to add chunk to buffer")
	}
	return packer.readBuffer(), nil
}

func (packer *Packer) Close() {
	C.ogg_opus_packer_destroy(packer.c_object)
}

func (packer *Packer) readBuffer() []byte {
	n := C.size_t(0)
	buffer := C.ogg_opus_paker_get_buffer(packer.c_object, &n)
	C.ogg_opus_packer_clear_buffer(packer.c_object)
	return C.GoBytes(unsafe.Pointer(buffer), C.int(n))
}

func opusDecoderStatusToError(status int) error {
	switch status {
	case C.OPUS_BAD_ARG:
		return errors.New("One or more invalid/out of range arguments")
	case C.OPUS_BUFFER_TOO_SMALL:
		return errors.New("The mode struct passed is invalid")
	case C.OPUS_INTERNAL_ERROR:
		return errors.New("An internal error was detected while docoding opus chunk")
	case C.OPUS_INVALID_PACKET:
		return errors.New("The compressed data passed is corrupted")
	case C.OPUS_UNIMPLEMENTED:
		return errors.New("Invalid/unsupported request number")
	case C.OPUS_INVALID_STATE:
		return errors.New(" An decoder structure is invalid or already freed")
	case C.OPUS_ALLOC_FAIL:
		return errors.New("Memory allocation has failed")
	default:
		return errors.New("Unexpected error")
	}
}

func initStatusToError(status int) error {
	switch status {
	case C.OGG_OPUS_PACKER_INIT_STATUS_OK:
		return nil
	case C.OGG_OPUS_PACKER_INIT_STATUS_STREAM_INIT_ERROR:
		return errors.New("Failed to init ogg stream")
	case C.OGG_OPUS_PACKER_INIT_STATUS_HEADER_ERROR:
		return errors.New("Failed to add a header to stream")
	case C.OGG_OPUS_PACKER_INIT_STATUS_ADD_TO_BUFFER_ERROR:
		return errors.New("Failed to add packet to buffer")
	case C.OGG_OPUS_PACKER_INIT_STATUS_COMMENT_ERROR:
		return errors.New("Failed to add comments to stream")
	default:
		return opusDecoderStatusToError(status)
	}
}
