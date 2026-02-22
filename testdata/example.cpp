// CodeOwner: @cpp_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in C++.
 */

#include <iostream>
#include <string>

int x = 42;  // Inline comment

std::string greet(const std::string &name) {
    // Another single-line comment
    return "Hello, " + name + "!";
}

int main() {
    std::cout << greet("world") << std::endl;  /* Inline block comment */
    return 0;
}
