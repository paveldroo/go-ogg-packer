#pragma once
#include <stdlib.h>
#include <ogg/ogg.h>

typedef struct buffer_t {
    void *data;
    size_t len;
    size_t read_index;
    size_t alloc;
} buffer_t;

buffer_t* buffer_create();

void buffer_destroy(buffer_t *s);

int buffer_add(buffer_t *s, const ogg_page *page);

char* buffer_get(buffer_t* s, size_t* n);

void buffer_reset(buffer_t *s);