#include <endian.h>
#include <stdlib.h>
#include <string.h>
#include "ogg_buffer.h"
#include "ogg_opus_packer.h"

#define MAX_FRAME_SIZE 5760
#define DEFAULT_SAMPLE_RATE 48000

ogg_opus_packer_t* ogg_opus_packer_create() {
    ogg_opus_packer_t* s = malloc(sizeof(ogg_opus_packer_t));
    if (!s) return NULL;
    if(!(s->buffer = buffer_create())) {
        free(s);
        return NULL;
    }
    return s;
}

static inline void fill_header(uint8_t *data, uint8_t channel_count, uint32_t sample_rate) {
    memcpy(data, "OpusHead", 8);

    data[8] = 1; /* Version number */
    data[9] = channel_count; /* Channels */

    *((uint16_t*) (data + 10)) = 0;
    *((uint32_t*) (data + 12)) = htole32(sample_rate);
    *((uint16_t*) (data + 16)) = 0;

    data[18] = 0; /* Mapping */
}

static int ogg_opus_packer_add_header(ogg_opus_packer_t *s) {
    uint8_t data[19];
    fill_header(data, s->channel_count, s->sample_rate);

    ogg_packet packet = {
    .packet = data,
    .bytes = 19,
    .b_o_s = 1,
    .e_o_s = 0,
    .granulepos = 0,
    .packetno = s->packetno++,
    };
    return ogg_stream_packetin(&s->stream_state, &packet);
}

static int ogg_opus_packer_add_comments(ogg_opus_packer_t *s) {
    uint8_t data[9];
    memcpy(data, "OpusTags", 9);
    ogg_packet packet = {
    .packet = data,
    .bytes = 9,
    .b_o_s = 0,
    .e_o_s = 0,
    .granulepos = 0,
    .packetno = s->packetno++,
    };
    return ogg_stream_packetin(&s->stream_state, &packet);
}

int ogg_opus_packer_init(ogg_opus_packer_t *s, uint8_t channel_count, uint32_t sample_rate, int serialno) {
    s->channel_count = channel_count;
    s->sample_rate = sample_rate;
    s->packetno = 1;
    s->granulepos = 0;
    int error;
    if (ogg_stream_init(&s->stream_state, serialno))
        return OGG_OPUS_PACKER_INIT_STATUS_STREAM_INIT_ERROR;
    if (!(s->opus_decoder = opus_decoder_create(DEFAULT_SAMPLE_RATE, channel_count, &error)))
        return error;
    if (ogg_opus_packer_add_header(s)) {
        return OGG_OPUS_PACKER_INIT_STATUS_HEADER_ERROR;
    }
    ogg_page page;
    while (ogg_stream_flush(&s->stream_state, &page)) {
        if (buffer_add(s->buffer, &page))
            return OGG_OPUS_PACKER_INIT_STATUS_ADD_TO_BUFFER_ERROR;
    }
    if (ogg_opus_packer_add_comments(s)) {
        return OGG_OPUS_PACKER_INIT_STATUS_COMMENT_ERROR;
    }
    while (ogg_stream_flush(&s->stream_state, &page)) {
        if (buffer_add(s->buffer, &page))
            return OGG_OPUS_PACKER_INIT_STATUS_ADD_TO_BUFFER_ERROR;
    }
    return OGG_OPUS_PACKER_INIT_STATUS_OK;
}

int ogg_opus_packer_add_opus_chunk(ogg_opus_packer_t *s, const void *data, size_t len, int eos, int samples_count) {
    int16_t buffer[MAX_FRAME_SIZE*s->channel_count*sizeof(opus_int16)];
    int num_samples_per_channel;
    if (samples_count < 0) {
        num_samples_per_channel = opus_decode(s->opus_decoder, data, len, buffer, MAX_FRAME_SIZE, 0);
        if (num_samples_per_channel < 0) {
            return num_samples_per_channel - 1; //Subtract one to avoid collision of errors with ogg_stream_packetin
        }
    } else {
        num_samples_per_channel = (samples_count * DEFAULT_SAMPLE_RATE) / (s->sample_rate * s->channel_count);
    }
    ogg_packet packet = {
    .packet = (void*) data,
    .bytes = len,
    .b_o_s = 0,
    .e_o_s = eos,
    .granulepos = s->granulepos,
    .packetno = s->packetno++,
    };
    s->granulepos += num_samples_per_channel;
    return ogg_stream_packetin(&s->stream_state, &packet);
}

int ogg_opus_packer_collect_pages(ogg_opus_packer_t *s) {
    ogg_page page;
    int status;
    while (ogg_stream_pageout(&s->stream_state, &page)) {
        if (status = buffer_add(s->buffer, &page))
            return status;
    }
    return 0;
}

int ogg_opus_packer_flush_pages(ogg_opus_packer_t *s) {
    ogg_page page;
    int status;
    while (ogg_stream_flush(&s->stream_state, &page)) {
        if (status = buffer_add(s->buffer, &page))
            return status;
	}
    return 0;
}

char* ogg_opus_paker_get_buffer(ogg_opus_packer_t *s, size_t *n) {
    return buffer_get(s->buffer, n);
}

void ogg_opus_packer_clear_buffer(ogg_opus_packer_t *s) {
    buffer_reset(s->buffer);
}

void ogg_opus_packer_destroy(ogg_opus_packer_t *s) {
    if (!s) return;
    if (s->opus_decoder)
        opus_decoder_destroy(s->opus_decoder);
    if (&s->stream_state)
        ogg_stream_clear(&s->stream_state);
    if (s->buffer)
        buffer_reset(s->buffer);
    if (s->buffer)
        buffer_destroy(s->buffer);
    free(s);
}
