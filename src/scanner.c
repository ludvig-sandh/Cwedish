#include "scanner.h"

#include <stdlib.h>
#include <stdbool.h>

char *read_full_file_content(const char *path, size_t *out_size) {
    FILE *f = fopen(path, "rb");
    if (!f) {
        return NULL;
    }

    if (fseek(f, 0, SEEK_END) != 0) {
        fclose(f);
        return NULL;
    }

    long len = ftell(f);
    if (len < 0) {
        fclose(f);
        return NULL;
    }

    rewind(f);

    char *buf = malloc((size_t)len);
    if (!buf) {
        fclose(f);
        return NULL;
    }

    size_t n = fread(buf, 1, (size_t)len, f);
    fclose(f);

    if (n != (size_t)len) {
        free(buf);
        return NULL;
    }

    if (out_size) {
        *out_size = n;
    }
    
    return buf;
}

TokenArray scan_content(const char *content, size_t size) {
    TokenArray token_array;
    init_token_array(&token_array);

    Token token;
    token.type = TOK_KEYWORD;
    token.start = content;
    token.length = 0;
    
    // TODO: Actually tokenize. This placeholder turns each word into a token.
    for (const char *c = content; c < content + size; ++c) {
        if (*c == ' ' || *c == '\r' || *c == '\n') {
            if (token.length != 0) {
                append_token(&token_array, &token);
            }

            token.start = c + 1;
            token.length = 0;
        }else {
            token.length++;
        }
    }

    if (token.length != 0) {
        append_token(&token_array, &token);
    }

    return token_array;
}
