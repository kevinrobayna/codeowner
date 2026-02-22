#!/usr/bin/perl

# CodeOwner: @perl_owner

# This is a single-line comment

=pod
This is a multi-line
block comment in Perl.
=cut

my $x = 42;  # Inline comment

sub greet {
    # Another single-line comment
    my ($name) = @_;
    return "Hello, $name!";
}

print greet("world") . "\n";
