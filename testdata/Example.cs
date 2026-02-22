// CodeOwner: @csharp_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in C#.
 */

/// <summary>XML doc comment.</summary>
class Example
{
    int x = 42;  // Inline comment

    string Greet(string name)
    {
        // Another single-line comment
        return $"Hello, {name}!";
    }

    static void Main(string[] args)
    {
        var ex = new Example();
        System.Console.WriteLine(ex.Greet("world"));  /* Inline block comment */
    }
}
