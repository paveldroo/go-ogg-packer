package opus_decoder

/*
#cgo pkg-config: opus
#cgo CFLAGS: -Wimplicit-function-declaration -Wall

#include <string.h>
#include <stdlib.h>
#include <opus.h>


typedef struct opus_decoder_wrapper_t {
  int num_channels;
  OpusDecoder* decoder;
  void *buffer;
} opus_decoder_wrapper_t;

#define MAX_FRAME_SIZE 5760


opus_decoder_wrapper_t* opus_decoder_wrapper_create(
  int num_channels,
  int sample_rate) {
  int error;
  OpusDecoder* decoder = opus_decoder_create(sample_rate, num_channels, &error);
  if (error != OPUS_OK)
    return NULL;

  opus_decoder_wrapper_t* s = malloc(sizeof(opus_decoder_wrapper_t) + MAX_FRAME_SIZE*num_channels*sizeof(float));
  if (!s) {
    if (decoder)
      opus_decoder_destroy(decoder);
    return NULL;
  }
  s->num_channels = num_channels;
  s->buffer = (void*) (s + 1);
  s->decoder = decoder;
  return s;
}
void opus_decoder_wrapper_destroy(opus_decoder_wrapper_t* s) {
  if (!s) return;
  if (s->decoder)
    opus_decoder_destroy(s->decoder);
  free(s);
}
int16_t* opus_decoder_wrapper_decode_int16(opus_decoder_wrapper_t* s, const uint8_t* bytes, int byte_count,
					   int *decoded_sample_count_or_errno) {
  int num_samples_per_channel = opus_decode(s->decoder, bytes, byte_count, s->buffer, MAX_FRAME_SIZE, 0);
  if (num_samples_per_channel < 0) {
    *decoded_sample_count_or_errno = num_samples_per_channel;
    return NULL;
  }
  *decoded_sample_count_or_errno = num_samples_per_channel*s->num_channels;
  return s->buffer;
}
float* opus_decoder_wrapper_decode_float(opus_decoder_wrapper_t* s, const uint8_t* bytes, int byte_count,
					 int *decoded_sample_count_or_errno) {
  int num_samples_per_channel = opus_decode_float(s->decoder, bytes, byte_count, s->buffer, MAX_FRAME_SIZE, 0);
  if (num_samples_per_channel < 0) {
    *decoded_sample_count_or_errno = num_samples_per_channel;
    return NULL;
  }
  *decoded_sample_count_or_errno = num_samples_per_channel*s->num_channels;
  return s->buffer;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type Decoder struct {
	c_object *C.opus_decoder_wrapper_t
}

func copySliceInt16(input []int16) []int16 {
	output := make([]int16, len(input))
	copy(output, input)
	return output
}
func copySliceFloat32(input []float32) []float32 {
	output := make([]float32, len(input))
	copy(output, input)
	return output
}

func New(num_channels int, sample_rate int) *Decoder {
	this := &Decoder{C.opus_decoder_wrapper_create(C.int(num_channels), C.int(sample_rate))}
	if this.c_object == nil {
		return nil
	}
	return this
}
func (this *Decoder) DecodeInt16(data []byte) ([]int16, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data block")
	}
	decoded_sample_count_or_errno := C.int(0)
	ret_samples := C.opus_decoder_wrapper_decode_int16(this.c_object, (*C.uint8_t)(&data[0]), C.int(len(data)),
		(&decoded_sample_count_or_errno))
	if ret_samples == nil {
		return nil, errors.New(C.GoString(C.opus_strerror(decoded_sample_count_or_errno)))
	}
	return copySliceInt16((*[1 << 28]int16)(unsafe.Pointer(ret_samples))[:int(decoded_sample_count_or_errno):int(decoded_sample_count_or_errno)]), nil
}
func (this *Decoder) DecodeFloat32(data []byte) ([]float32, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data block")
	}
	decoded_sample_count_or_errno := C.int(0)
	ret_samples := C.opus_decoder_wrapper_decode_float(this.c_object, (*C.uint8_t)(&data[0]), C.int(len(data)),
		(&decoded_sample_count_or_errno))
	if ret_samples == nil {
		return nil, errors.New(C.GoString(C.opus_strerror(decoded_sample_count_or_errno)))
	}
	return copySliceFloat32((*[1 << 28]float32)(unsafe.Pointer(ret_samples))[:int(decoded_sample_count_or_errno):int(decoded_sample_count_or_errno)]), nil
}
func (this *Decoder) Close() {
	C.opus_decoder_wrapper_destroy(this.c_object)
}
