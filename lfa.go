package main

import (
	"fmt"
	"os"
	"regex"
)

func main() {
	var word string
	var re string
	fmt.Scanf("%s", &re)
	fmt.Scanf("%s", &word)
	nfa := regex.RegexToNFA(re)
	nfa.ToDFA()
	nfa.Minimize()
	fmt.Printf("%v", nfa.Check(word))
	fmt.Printf("\n")
	nfa.Print(os.Stdout)
}
