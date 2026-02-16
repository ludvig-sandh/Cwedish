CC = gcc
CFLAGS = -std=c17 -Wall -Wextra -pedantic -O2

TARGET = cwedish
SRC = main.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(SRC) -o $(TARGET)

clean:
	rm -f $(TARGET)

.PHONY: all clean
