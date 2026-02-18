#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#include "mappings_parser.h"

int main() {
    const char *mappings_file_path = "keywords-table.txt";

    Mapping mappings[MAX_NUM_MAPPINGS];
    size_t num_mappings = parse_mappings(mappings_file_path, mappings);
    for (size_t i = 0; i < num_mappings; i++) {
        Mapping *mapping = &mappings[i];
        printf("Parsed keyword \"%s\" -> \"%s\"\n", mapping->original, mapping->translated);
    }
}
