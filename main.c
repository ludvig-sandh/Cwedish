#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#define MAX_KEYWORD_LENGTH 15
#define MAX_NUM_MAPPINGS 40

typedef struct {
    char original[MAX_KEYWORD_LENGTH + 1];
    char translated[MAX_KEYWORD_LENGTH + 1];
} Mapping;

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

int main() {
    const char *mappings_file_path = "keywords-table.txt";

    Mapping mappings[MAX_NUM_MAPPINGS];
    size_t num_mappings = parse_mappings(mappings_file_path, mappings);
    for (size_t i = 0; i < num_mappings; i++) {
        Mapping *mapping = &mappings[i];
        printf("Parsed keyword \"%s\" -> \"%s\"\n", mapping->original, mapping->translated);
    }
}
