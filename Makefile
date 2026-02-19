CC = gcc
CFLAGS = -std=c17 -Wall -Wextra -pedantic -O2
DEBUG_CFLAGS = -std=c17 -Wall -Wextra -pedantic -g -O0

TARGET = cwedish
DEBUG_TARGET = cwedish-debug
SRC = src/main.c src/mappings_parser.c src/token.c src/scanner.c

all: $(TARGET)

debug: $(DEBUG_TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(SRC) -o $(TARGET)

$(DEBUG_TARGET): $(SRC)
	$(CC) $(DEBUG_CFLAGS) $(SRC) -o $(DEBUG_TARGET)

clean:
	rm -f $(TARGET) $(DEBUG_TARGET)

.PHONY: all debug clean
