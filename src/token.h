#ifndef TOKEN_H
#define TOKEN_H

#include <stdlib.h>
#include <assert.h>

typedef enum {
    TOK_KEYWORD,
    TOK_IDENTIFIER,
    TOK_OPERATOR,
    TOK_LITERAL,
    TOK_EOF
    // Something like this. TODO: Find out actual needed types
} TokenType;

typedef struct {
    TokenType type;
    const char *start;
    size_t length;
} Token;

// Dynamic array
typedef struct {
    size_t capacity;
    size_t size;
    Token *data;
} TokenArray;

void init_token_array(TokenArray *array);

void free_token_array(TokenArray *array);

void append_token(TokenArray *array, Token token);

Token get_token(TokenArray *array, size_t index);

#endif