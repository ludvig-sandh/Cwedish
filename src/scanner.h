#ifndef SCANNER_H
#define SCANNER_H

#include "token.h"

#include <stdio.h>

char *read_full_file_content(const char *path, size_t *out_size);

TokenArray scan_content(const char *content, size_t size);

#endif
