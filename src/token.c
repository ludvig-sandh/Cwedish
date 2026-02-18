#include "token.h"

void init_token_array(TokenArray *array) {
    array->capacity = 0;
    array->size = 0;
    array->data = NULL;
}

void free_token_array(TokenArray *array) {
    free(array->data);
    array->data = NULL;
    array->capacity = 0;
    array->size = 0;
}

void grow_token_array(TokenArray *array) {
    size_t new_cap = array->capacity == 0 ? 1 : array->capacity * 2;

    Token *new_data = realloc(array->data, new_cap * sizeof(Token));
    assert(new_data != NULL && "realloc failed");

    array->data = new_data;
    array->capacity = new_cap;
}

void append_token(TokenArray *array, const Token *token) {
    if (array->size == array->capacity) {
        grow_token_array(array);
    }

    array->data[array->size++] = *token;
}

Token *get_token(TokenArray *array, size_t index) {
    assert(index < array->size && "index out of bounds for get_token()");
    return &array->data[index];
}
