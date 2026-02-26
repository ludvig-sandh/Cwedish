#ifndef MAPPINGS_PARSER_H
#define MAPPINGS_PARSER_H

#include "token.h"

#include <stdio.h>
#include <stdbool.h>

#define MAX_KEYWORD_LENGTH 15
#define MAX_NUM_MAPPINGS 40

typedef struct {
    char original[MAX_KEYWORD_LENGTH + 1];
    char translated[MAX_KEYWORD_LENGTH + 1];
} Mapping;

size_t parse_mappings(const char *path, Mapping *mappings);

bool translate(TokenArray *array, Mapping *mappings, size_t num_mappings, const char *path);

#endif
