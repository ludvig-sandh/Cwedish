#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#include "mappings_parser.h"
#include "token.h"
#include "scanner.h"

int main() {
    const char *mappings_file_path = "keywords-table.txt";

    Mapping mappings[MAX_NUM_MAPPINGS];
    size_t num_mappings = parse_mappings(mappings_file_path, mappings);
    // for (size_t i = 0; i < num_mappings; i++) {
    //     Mapping *mapping = &mappings[i];
    //     printf("Parsed keyword \"%s\" -> \"%s\"\n", mapping->original, mapping->translated);
    // }

    size_t out_size = 0;
    const char *content = read_full_file_content("example.c", &out_size);
    assert(content != NULL && "couldn't read file contents");

    TokenArray token_array = scan_content(content, out_size);

    for (size_t i = 0; i < token_array.size; ++i) {
        Token *token = get_token(&token_array, i);
        printf("Token %d: \'", (int)(i + 1));
        for (const char *c = token->start; c < token->start + token->length; ++c) {
            printf("%c", *c);
        }
        printf("\'\n");
    }
}
