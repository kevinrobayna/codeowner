// CodeOwner: @scala_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in Scala.
 */

object Example {

  val x: Int = 42  // Inline comment

  def greet(name: String): String = {
    // Another single-line comment
    s"Hello, $name!"
  }

  def main(args: Array[String]): Unit = {
    println(greet("world"))  /* Inline block comment */
  }
}
