#pragma once
#include <opus/opus.h>
#include <stdint.h>
#include <ogg/ogg.h>
#include "ogg_buffer.h"

typedef enum ogg_opus_packer_init_status {
    OGG_OPUS_PACKER_INIT_STATUS_OK = 0,
    OGG_OPUS_PACKER_INIT_STATUS_STREAM_INIT_ERROR = -100,
    OGG_OPUS_PACKER_INIT_STATUS_HEADER_ERROR,
    OGG_OPUS_PACKER_INIT_STATUS_ADD_TO_BUFFER_ERROR,
    OGG_OPUS_PACKER_INIT_STATUS_COMMENT_ERROR,
} packer_init_status;

typedef struct ogg_opus_packer_t {
    uint8_t channel_count;
    uint32_t sample_rate;
    ogg_int64_t packetno;
    ogg_int64_t granulepos;
    ogg_stream_state stream_state;
    buffer_t *buffer;
    OpusDecoder *opus_decoder;
} ogg_opus_packer_t;

ogg_opus_packer_t* ogg_opus_packer_create();

int ogg_opus_packer_init(ogg_opus_packer_t *s, uint8_t channel_count, uint32_t sample_rate, int serialno);

void ogg_opus_packer_uninit(ogg_opus_packer_t *s);

void ogg_opus_packer_destroy(ogg_opus_packer_t *s);

int ogg_opus_packer_add_opus_chunk(ogg_opus_packer_t *s, const void *data, size_t len, int eos, int samples_count);

int ogg_opus_packer_collect_pages(ogg_opus_packer_t *s);

int ogg_opus_packer_flush_pages(ogg_opus_packer_t *s);

char* ogg_opus_paker_get_buffer(ogg_opus_packer_t *s, size_t *n);

void ogg_opus_packer_clear_buffer(ogg_opus_packer_t *s);