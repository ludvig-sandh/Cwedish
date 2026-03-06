# Cwedish
A translator for a custom C language where all the keywords are replaced by Swedish words

Cwedish example:
```c
hel main() {
    för (hel i = 0; i < 3; ++i) {
        om (i == 1) {
            fortsätt;
        }
    }
    returnera 0;
}
```
Running the translator on this would yield regular valid C:
```c
int main() {
    for (int i = 0; i < 3; ++i) {
        if (i == 1) {
            continue;
        }
    }
    return 0;
}
```

Since Cwedish allows you to provide your own set of keyword translations with a custom dictionary, you could create your own version of C in another language, for example Cpanish, Corean, GreeC, Latin-C, Corwegian or C-hinese.

## Build
Requires Go to be installed.

```sh
make build
```

By default, the `dictionary.txt` file beside the executable is used for translating keywords, but a custom dictionary can instead be passed with `-d`.

## Usage
Translate main.cwe into main.c:
```sh
./cwedish main.cwe
```
Translate input.cwe into output.c
```sh
./cwedish -o output.c input.cwe
```
Translate main.cwe into main.c using a custom keyword translation dictionary
```sh
./cwedish -d ./custom-dictionary.txt main.cwe
```

Flags:
- `-o` sets the output file path. If omitted, the output file gets the same name as the input file with a `.c` extension.
- `-d` sets a custom dictionary file path. If omitted, the program looks for `dictionary.txt` beside the executable.

## Tests
```sh
make test
```
