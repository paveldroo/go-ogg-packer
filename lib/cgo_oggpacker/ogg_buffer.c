#include <string.h>
#include "ogg_buffer.h"

#define ALLOCATION_SIZE 4096

buffer_t* buffer_create() {
    buffer_t* s = malloc(sizeof(buffer_t));
    if (!s) return NULL;
    void* data = malloc(ALLOCATION_SIZE);
    if (!data) {
        free(s);
        return NULL;
    }
    s->data = data;
    s->len = 0;
    s->read_index = 0;
    s->alloc = ALLOCATION_SIZE;
    return s;
}

void buffer_destroy(buffer_t *s) {
    if (!s) return;
    free(s->data);
    free(s);
}

int buffer_add(buffer_t *s, const ogg_page *page) {
    size_t len = s->len;
    size_t header_len = page->header_len;
    size_t body_len = page->body_len;
    size_t fit_alloc = len + header_len + body_len;
    size_t alloc = s->alloc;
    void *data = s->data;
    if (alloc < fit_alloc) {
        do {
            alloc <<= 1;
        } while (alloc < fit_alloc);
        if (!(data = realloc(data, alloc)))
            return -1;
        s->alloc = alloc;
        s->data = data;
    }
    memcpy(data + len, page->header, header_len);
    len += header_len;
    memcpy(data + len, page->body, body_len);
    len += body_len;
    s->len = len;
    return 0;
}

char* buffer_get(buffer_t* s, size_t* n) {
    *n = s->len - s->read_index;
    return s->data + s->read_index;
}

void buffer_reset(buffer_t* s) {
    s->read_index = 0;
    s->len = 0;
}
