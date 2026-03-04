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

// Finite state machine for the different scanning modes
typedef void (*State)(TokenArray *, Token *, const char *);
State state;

void state_possibly_multi_char_operator(TokenArray *array, Token *token, const char *c);
void state_single_quote_string(TokenArray *array, Token *token, const char *c);
void state_double_quote_string(TokenArray *array, Token *token, const char *c);

void state_regular_code(TokenArray *array, Token *token, const char *c) {
    switch (*c) {
    case '\n':
    case '\r':
    case '\t':
    case ' ':
        // All these work like a separator I think
        append_and_reset_token(array, token);
        token->start = c;
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        break;
    case '\'':
        append_and_reset_token(array, token);
        token->start = c;
        token->length++;
        state = state_single_quote_string;
        break;
    case '\"':
        append_and_reset_token(array, token);
        token->start = c;
        token->length++;
        state = state_double_quote_string;
        break;
    case '{':
    case '}':
    case '(':
    case ')':
    case '=':
    case ',':
    case ';':
        // These should be tokens alone
        append_and_reset_token(array, token);
        token->start = c;
        token->length++;
        append_and_reset_token(array, token); // Apend this char as another token
        token->start = c + 1;
        break;
    case '+':
    case '-':
    case '/':
    case '*':
    case '|':
    case '&':
    case '^':
    case '~':
        // Depending on the next character (eg. '='), the token may need the next char as well
        state = state_possibly_multi_char_operator;
        token->length++;
        break;
    default:
        token->length++;
        break;
    }
}

void state_escaped_single_quote_string(TokenArray *array, Token *token, const char *c) {
    (void)array;
    (void)c;
    token->length++;
    state = state_single_quote_string;
}

void state_escaped_double_quote_string(TokenArray *array, Token *token, const char *c) {
    (void)array;
    (void)c;
    token->length++;
    state = state_double_quote_string;
}

void state_single_quote_string(TokenArray *array, Token *token, const char *c) {
    switch (*c) {
    case '\\':
        token->length++;
        state = state_escaped_single_quote_string;
        break;
    case '\'':
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
        break;
    default:
        token->length++;
        break;
    }
}

void state_double_quote_string(TokenArray *array, Token *token, const char *c) {
    switch (*c) {
    case '\\':
        token->length++;
        state = state_escaped_double_quote_string;
        break;
    case '\"':
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
        break;
    default:
        token->length++;
        break;
    }
}

void state_single_line_comment(TokenArray *array, Token *token, const char *c) {
    if (*c == '\r' || *c == '\n') {
        append_and_reset_token(array, token);
        token->start = c;
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
    }else {
        token->length++;
    }
}

void state_multi_line_comment(TokenArray *array, Token *token, const char *c) {
    token->length++;
    if (*(c - 1) == '*' && *c == '/') {
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
    }
}

void state_possibly_multi_char_operator(TokenArray *array, Token *token, const char *c) {
    if (*c == '=') { // +=, |= or similar
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
    }else if (*(c - 1) == '/' && *c == '/') { // //
        token->length++;
        state = state_single_line_comment;
    }else if (*(c - 1) == '/' && *c == '*') { // /*
        token->length++;
        state = state_multi_line_comment;
    }else if (*(c - 1) == '+' && *c == '+') { // ++
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
    }else if (*(c - 1) == '-' && *c == '-') { // --
        token->length++;
        append_and_reset_token(array, token);
        token->start = c + 1;
        state = state_regular_code;
    }else if (*(c - 1) == '<' && *c == '<') { // <<
        token->length++;
        append_and_reset_token(array, token);
        state = state_possibly_multi_char_operator;
    }else if (*(c - 1) == '>' && *c == '>') { // >>
        token->length++;
        append_and_reset_token(array, token);
        state = state_possibly_multi_char_operator;
    }else {
        append_and_reset_token(array, token);
        token->start = c;
        state = state_regular_code;
    }
}

TokenArray scan_content(const char *content, size_t size) {
    TokenArray token_array;
    init_token_array(&token_array);

    Token token;
    token.type = TOK_KEYWORD;
    token.start = content;
    token.length = 0;

    state = state_regular_code;
    
    for (const char *c = content; c < content + size; ++c) {
        state(&token_array, &token, c);
    }

    append_and_reset_token(&token_array, &token);

    return token_array;
}
