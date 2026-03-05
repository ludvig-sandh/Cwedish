#include <stdio.h>

typedef struct Demo {
    short small;
    int count;
    long big;
    char letter;
    float ratio;
    double precise;
} Demo;

static int choose_value(int input) {
    if (input < 0) {
        goto fallback;
    }

    switch (input) {
    case 0:
        return 10;
    case 1:
        return 20;
    default:
        return input;
    }

fallback:
    return -1;
}

int main() {
    Demo demo = {1, 2, 3, 'A', 1.5f, 2.5};
    int result = 0;

    for (int i = 0; i < 5; ++i) {
        if (i == 2) {
            continue;
        }

        result += choose_value(i);

        if (result > 20) {
            break;
        }
    }

    if (demo.count > 0) {
        result += demo.small;
    } else {
        result -= demo.small;
    }

    while (result > 100) {
        result--;
    }

    return result;
}
