# Cwedish
A translator for a custom C language where all the keywords are replaced by Swedish words

## Build
```sh
go build -o cwedish ./cmd
```

By default, the `dictionary.txt` file beside the executable is used for translating keywords, but a custom dictionary can instead be passed with `-d`.

## Usage
```sh
./cwedish main.cwe
./cwedish -o output.c input.cwe
./cwedish -d ./custom-dictionary.txt main.cwe
```

Flags:
- `-o` sets the output file path. If omitted, the output file gets the same name as the input file with a `.c` extension.
- `-d` sets a custom dictionary file path. If omitted, the program looks for `dictionary.txt` beside the executable.

## Tests
```sh
go test ./...
```
