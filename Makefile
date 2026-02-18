CC = gcc
CFLAGS = -std=c17 -Wall -Wextra -pedantic -O2

TARGET = cwedish
SRC = src/main.c src/mappings_parser.c src/token.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(SRC) -o $(TARGET)

clean:
	rm -f $(TARGET)

.PHONY: all clean
