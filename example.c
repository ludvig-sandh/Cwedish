#include <stdio.h>

int main() {
    char text[] = "för(;;) {}\n";
    // för
    for (int i = 0; i < 5; ++i) {
        /*A
        B*/
        printf("%d\n", i);
    }
}
