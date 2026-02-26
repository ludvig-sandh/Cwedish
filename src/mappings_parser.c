#include "mappings_parser.h"
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <math.h>

bool copy_until_stop_character(const char *line, char stop, char *dest, size_t *num_copied) {
    char *stop_pos = strchr(line, stop);
    if (stop_pos == NULL) {
        return true;
    }
    
    size_t index = stop_pos - line;
    if (index > MAX_KEYWORD_LENGTH) {
        return true;
    }

    memcpy(dest, line, index);
    dest[index] = '\0';

    *num_copied = index;
    return false;
}

bool parse_mapping_line(const char *line, Mapping *mapping) {
    size_t keyword_size;
    bool err = copy_until_stop_character(line, ':', mapping->original, &keyword_size);
    if (err) {
        printf("Couldn't parse the original keyword\n");
        return true;
    }

    line += keyword_size + 1;
    err = copy_until_stop_character(line, '\n', mapping->translated, &keyword_size);
    if (err) {
        printf("Couldn't parse the translated keyword\n");
        return true;
    }

    return false;
}

size_t parse_mappings(const char *path, Mapping *mappings) {
    FILE *fptr = fopen(path, "r");

    const int line_length = MAX_KEYWORD_LENGTH * 2 + 2; // colon and newline
    char line[line_length + 1];
    int lines_parsed = 0;
    while (fgets(line, sizeof(line), fptr) != NULL) {
        Mapping *mapping = &mappings[lines_parsed];
        if (lines_parsed == MAX_NUM_MAPPINGS) {
            printf("Warning: more than max supported mappings found. Only parsed the first %d.\n", MAX_NUM_MAPPINGS);
            break;
        }
        
        bool err = parse_mapping_line(line, mapping);
        if (err) {
            printf("Couldn't parse line %d. Exiting.\n", lines_parsed + 1);
            exit(1);
        }
        
        lines_parsed++;
    }
    
    return lines_parsed;
}

Mapping *get_mapping(Token *token, Mapping *mappings, size_t num_mappings) {
    for (size_t i = 0; i < num_mappings; ++i) {
        size_t safe_comparison_length = token->length < MAX_KEYWORD_LENGTH ? token->length : MAX_KEYWORD_LENGTH;
        if (memcmp(token->start, mappings[i].translated, safe_comparison_length) == 0) {
            return &mappings[i];
        }
    }
    
    return NULL;
}

bool translate(TokenArray *array, Mapping *mappings, size_t num_mappings, const char *path) {
    FILE *fptr = fopen(path, "wb");
    if (fptr == NULL) {
        perror("fopen");
        return false;
    }

    for (size_t i = 0; i < array->size; ++i) {
        Token *token = get_token(array, i);
        Mapping *mapping_found = get_mapping(token, mappings, num_mappings);
        if (mapping_found == NULL) {
            size_t written = fwrite(token->start, sizeof(char), token->length, fptr);
            if (written != token->length) {
                perror("fwrite");
                fclose(fptr);
                return false;
            }
        }else {
            int written = fprintf(fptr, "%s", mapping_found->original);
            if (written != (int)strlen(mapping_found->original)) {
                perror("fprintf");
                fclose(fptr);
                return false;
            }
        }
    }

    fclose(fptr);
    return true;
}
