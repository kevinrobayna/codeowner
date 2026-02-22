// CodeOwner: @kotlin_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in Kotlin.
 */

/** KDoc comment for the function. */
fun greet(name: String): String {
    // Another single-line comment
    return "Hello, $name!"
}

fun main() {
    val x = 42  // Inline comment
    val result = greet("world")  /* Inline block comment */
    println(result)
}
