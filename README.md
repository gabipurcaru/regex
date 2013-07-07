Regex
=====

This is a very simple regular expression engine, created for learning purposes. It supports parantheses, 
the [Kleene star](https://en.wikipedia.org/wiki/Kleene_star), and the OR operator.

It works by turning a regular expression into a [λ-NFA](http://en.wikipedia.org/wiki/Nondeterministic_finite_automaton_with_%CE%B5-moves), 
turns that into a simple [NFA](http://en.wikipedia.org/wiki/Nondeterministic_finite_automaton), 
then into a [DFA](http://en.wikipedia.org/wiki/Deterministic_finite_automaton), then 
[minimizes the DFA](http://en.wikipedia.org/wiki/DFA_minimization). A word is then matched by the 
regular expression if it's recognized by the resulting DFA.

The package is written in [golang](http://golang.org/) and it doesn't have any external dependencies.
Also, it doesn't rely on the `re` package.


What's this all about
================

If you want to learn about Finite Automata, this package should be hacking material.


How do I use this
==============

`go run lfa.go` -- it will prompt you for a regular expression and a word to match against that expression.

It then prints whether the word matches or not, and also prints out the minimized DFA as a directed graph with
characters on edges.

You can use UTF-8 for the word/regex, though don't use the lambda character.

Example regular expressions:

    a*b
    (a|b)*c
    (word)*|anotherword
