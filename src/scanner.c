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

void append_and_reset_token(TokenArray *array, Token *token) {
    if (token->length != 0) {
        append_token(array, token);
    }

    token->start = NULL;
    token->length = 0;
}

TokenArray scan_content(const char *content, size_t size) {
    TokenArray token_array;
    init_token_array(&token_array);

    Token token;
    token.type = TOK_KEYWORD;
    token.start = content;
    token.length = 0;
    
    for (const char *c = content; c < content + size; ++c) {
        switch (*c) {
        case '\n':
        case '\r':
        case '\t':
        case ' ':
            // All these work like a separator I think
            append_and_reset_token(&token_array, &token); // Append current token and skip this
            token.start = c + 1;
            break;
        case '{':
        case '}':
        case '(':
        case ')':
        case '=':
        case ',':
        case ';':
            // These should be tokens alone
            append_and_reset_token(&token_array, &token); // Append current token
            token.start = c;
            token.length++;
            append_and_reset_token(&token_array, &token); // Apennd this char as another token
            token.start = c + 1;
            break;
        default:
            token.length++;
            break;
        }
    }

    append_and_reset_token(&token_array, &token);

    return token_array;
}
