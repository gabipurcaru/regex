package regex

import (
	"fmt"
	"nfa"
    // "os"
)

func RegexToNFA(re string) nfa.NFA {
	if len(re) == 0 {
		fmt.Printf("EMPTY STRING\n")
		return nfa.New()
	}
	if re[0] == '(' && re[len(re)-1] == ')' {
		fmt.Printf("ENCLOSED BY PARENS: %s\n", re)
		return RegexToNFA(re[1 : len(re)-1])
	}
	first_star_pos := -1
	first_pipe_pos := -1
	first_star_nesting := 1 << 30
	first_pipe_nesting := 1 << 30
	nesting := 0
	for pos, char := range re {
		if char == '(' {
			nesting++
		} else if char == ')' {
			nesting--
		} else if char == '|' {
			if nesting < first_pipe_nesting {
				first_pipe_nesting = nesting
				first_pipe_pos = pos
			}
		} else if char == '*' {
			if nesting < first_star_nesting {
				first_star_nesting = nesting
				first_star_pos = pos
			}
		}
	}
	if first_pipe_pos != -1 && first_pipe_nesting == 0 {
		n1 := RegexToNFA(re[0:first_pipe_pos])
		n2 := RegexToNFA(re[first_pipe_pos+1:])
		fmt.Printf("PIPE: %s <--> %s\n", re[0:first_pipe_pos], re[first_pipe_pos+1:])
		return nfa.Either(n1, n2)
	} else if first_star_pos != -1 && first_star_nesting == 0 {
		i := 0
		if re[first_star_pos-1] == ')' {
			reverse_nesting := 1
			i = first_star_pos - 2
			for ; ; i-- {
				if re[i] == ')' {
					reverse_nesting++
				} else if re[i] == '(' {
					reverse_nesting--
					if reverse_nesting == 0 {
						break
					}
				}
			}
		} else {
			i = first_star_pos - 1
		}
		res := nfa.New()
		if i > 0 {
			n1 := RegexToNFA(re[0:i])
			n2 := RegexToNFA(re[i:first_star_pos])
			res = nfa.Concat(n1, nfa.Star(n2))
		} else {
			res = nfa.Star(RegexToNFA(re[i:first_star_pos]))
		}
		if len(re[first_star_pos+1:]) > 0 {
			res = nfa.Concat(res, RegexToNFA(re[first_star_pos+1:]))
		}
		return res
	} else {
		if re[0] == '(' {
			nesting = 1
			i := 1
			for ; ; i++ {
    				if re[i] == '(' {
					nesting++
				} else if re[i] == ')' {
					nesting--
					if nesting <= 0 {
						break
					}
				}
			}
			fmt.Printf("PAREN: %s <--> %s\n", re[1:i], re[i+1:])
			return nfa.Concat(RegexToNFA(re[1:i]), RegexToNFA(re[i+1:]))
		} else {
			res := nfa.New()
			node := 1
			res.EntryState = 1
			res.NumStates = 1
			res.FinalStates = []int{1}
			for pos, char := range re {
				if char == '(' {
					res.FinalStates = []int{node}
					fmt.Printf("PAREN2: %s <--> %s\n", re[:pos], re[pos:])
					return nfa.Concat(res, RegexToNFA(re[pos:]))
				}
				res.Graph[node] = map[rune][]int{char: []int{node + 1}}
				res.NumStates++
				res.NumTransitions++
				node++
			}
			res.FinalStates = []int{node}
			fmt.Printf("REGULAR: %s\n", re)
            // fmt.Printf("\n>>\n")
            // res.Print(os.Stdout)
            // fmt.Printf("\n<<\n")
			return res
		}
	}
	return nfa.New()
}
