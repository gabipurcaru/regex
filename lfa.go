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
	nfa.Print(os.Stdout)
}
