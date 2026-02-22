// CodeOwner: @java_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in Java.
 */

/** Javadoc comment for the class. */
public class Example {

    int x = 42;  // Inline comment

    public String greet(String name) {
        // Another single-line comment
        return "Hello, " + name + "!";
    }

    public static void main(String[] args) {
        Example ex = new Example();
        System.out.println(ex.greet("world"));  /* Inline block comment */
    }
}
