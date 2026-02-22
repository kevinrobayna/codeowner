// CodeOwner: @c_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in C.
 */

#include <stdio.h>

int x = 42;  // Inline comment

void greet(const char *name) {
    // Another single-line comment
    printf("Hello, %s!\n", name);
}

int main(void) {
    greet("world");  /* Inline block comment */
    return 0;
}
