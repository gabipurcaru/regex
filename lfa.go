package main

import (
	"fmt"
	"os"
	"regex"
)

func main() {
	var word string
	var re string
	fmt.Print("Give a regular expression: ")
	fmt.Scanf("%s", &re)
	fmt.Print("Give a word to match: ")
	fmt.Scanf("%s", &word)
	nfa := regex.RegexToNFA(re)
	dfa := nfa.ToDFA()
	dfa.Minimize()
	fmt.Printf("Matches? %v\n", dfa.Check(word))

	fmt.Print("Minimized DFA for the regular expression: \n")
	dfa.Print(os.Stdout)
}
