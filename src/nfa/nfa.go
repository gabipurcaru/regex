package nfa

import (
	"dfa"
	"queue"
)

type NFA struct {
	dfa.DFA
}

func New() NFA {
	return NFA{dfa.NewDFA()}
}

func Copy(n NFA) NFA {
	res := New()
	res.NumStates = n.NumStates
	res.NumTransitions = n.NumTransitions
	res.EntryState = n.EntryState
	res.FinalStates = make([]int, len(n.FinalStates))
	copy(res.FinalStates, n.FinalStates)
	for node, _ := range n.Graph {
		res.Graph[node] = make(map[rune][]int)
		for character, _ := range n.Graph[node] {
			res.Graph[node][character] = make([]int, 0, len(n.Graph[node][character]))
			for _, neighbour := range n.Graph[node][character] {
				res.Graph[node][character] = append(res.Graph[node][character], neighbour)
			}
		}
	}
	return res
}

// Concatenates two NFAs. If n1 matched φ and n2 matched ψ, the resulting
// NFA will match φψ
func Concat(n1 NFA, n2 NFA) (n3 NFA) {
	n1 = Copy(n1)
	n2 = Copy(n2)
	n3 = New()
	offset := n1.NumStates
	n3.NumStates = n1.NumStates + n2.NumStates + 1
	n3.NumTransitions = n1.NumTransitions + n2.NumTransitions + len(n1.FinalStates) + 1
	n3.EntryState = n1.EntryState
	n3.FinalStates = make([]int, 0, len(n2.FinalStates))
	for _, node := range n2.FinalStates {
		n3.FinalStates = append(n3.FinalStates, node+offset)
	}
	n3.Graph = n1.Graph
	for node, _ := range n2.Graph {
		for character, _ := range n2.Graph[node] {
			for _, neighbour := range n2.Graph[node][character] {
				if _, ok := n3.Graph[node+offset]; !ok {
					n3.Graph[node+offset] = make(map[rune][]int)
				}
				if _, ok := n3.Graph[node+offset][character]; !ok {
					n3.Graph[node+offset][character] = make([]int, 0)
				}
				n3.Graph[node+offset][character] = append(
					n3.Graph[node+offset][character], neighbour+offset)
			}
		}
	}
	for _, node := range n1.FinalStates {
		if _, ok := n3.Graph[node]; !ok {
			n3.Graph[node] = make(map[rune][]int)
		}
		if _, ok := n3.Graph[node]['λ']; !ok {
			n3.Graph[node]['λ'] = make([]int, 0, 1)
		}
		n3.Graph[node]['λ'] = append(n3.Graph[node]['λ'], n3.NumStates)
	}
	n3.Graph[n3.NumStates] = map[rune][]int{'λ': []int{n2.EntryState + offset}}
	return
}

// OR-ing two NFAs together. If n1 matched φ and n2 matched ψ, the resulting
// NFA will match either of φ and ψ
func Either(n1 NFA, n2 NFA) (n3 NFA) {
	n1 = Copy(n1)
	n2 = Copy(n2)
	n3 = New()
	n3.NumStates = n1.NumStates + n2.NumStates + 1
	n3.NumTransitions = n1.NumTransitions + n2.NumTransitions + 2
	n3.EntryState = n3.NumStates // the last, newly added state
	n3.FinalStates = make([]int, len(n1.FinalStates), len(n1.FinalStates)+len(n2.FinalStates))
	offset := n1.NumStates
	copy(n3.FinalStates, n1.FinalStates)
	for _, node := range n2.FinalStates {
		n3.FinalStates = append(n3.FinalStates, node+offset)
	}
	n3.Graph = n1.Graph
	for node, _ := range n2.Graph {
		for character, _ := range n2.Graph[node] {
			for _, neighbour := range n2.Graph[node][character] {
				if _, ok := n3.Graph[node+offset]; !ok {
					n3.Graph[node+offset] = make(map[rune][]int)
				}
				if _, ok := n3.Graph[node+offset][character]; !ok {
					n3.Graph[node+offset][character] = make([]int, 0, 1)
				}
				n3.Graph[node+offset][character] = append(
					n3.Graph[node+offset][character], neighbour+offset)
			}
		}
	}
	n3.Graph[n3.EntryState] = map[rune][]int{
		'λ': []int{n1.EntryState, n2.EntryState + offset},
	}
	return
}

// Star operation on an NFA. If the original NFA matches φ, the resulting
// NFA will match zero or more occurrences of φ
func Star(n1 NFA) (n2 NFA) {
	n1 = Copy(n1)
	n2 = New()
	n2.NumStates = n1.NumStates + 1
	n2.NumTransitions = n1.NumTransitions + len(n1.FinalStates) + 1
	n2.EntryState = n1.EntryState
	n2.FinalStates = []int{n1.EntryState}
	n2.Graph = n1.Graph
	for _, state := range n1.FinalStates {
		if _, ok := n2.Graph[state]; !ok {
			n2.Graph[state] = make(map[rune][]int)
		}
		if _, ok := n2.Graph[state]['λ']; !ok {
			n2.Graph[state]['λ'] = make([]int, 0, 1)
		}
		n2.Graph[state]['λ'] = append(n2.Graph[state]['λ'], n2.NumStates)
	}
	n2.Graph[n2.NumStates] = map[rune][]int{'λ': []int{n1.EntryState}}
	return
}

// Transforms an NFA into a DFA that accepts the same language; some
// inaccessible states will be lost in the process
func (n *NFA) ToDFA() dfa.DFA {
	// first follow along λ-transitions to make them obsolete
	is_final := make([]bool, n.NumStates+1)
	for _, node := range n.FinalStates {
		is_final[node] = true
	}
	for i := 1; i <= n.NumStates; i++ {
		q := queue.New(n.NumStates)
		added := make([]bool, n.NumStates+1)
		for _, node := range n.Graph[i]['λ'] {
			q.Push(node)
			added[node] = true
		}
		for !q.Empty() {
			node, _ := q.Pop()
			if is_final[node] && !is_final[i] {
				is_final[i] = true
				n.FinalStates = append(n.FinalStates, i)
			}
			for character, neighbours := range n.Graph[node] {
				for _, neighbour := range neighbours {
					if character == 'λ' && !added[neighbour] {
						q.Push(neighbour)
					}
					n.Graph[i][character] = append(n.Graph[i][character], neighbour)
					n.NumTransitions++
				}
			}
		}
	}

	// then remove λ-transitions
	for i := 1; i <= n.NumStates; i++ {
		_, ok := n.Graph[i]['λ']
		if ok {
			n.NumTransitions -= len(n.Graph[i]['λ'])
			delete(n.Graph[i], 'λ')
		}
	}

	// then turn that into a DFA
	visited := make([]bool, 1<<uint(n.NumStates+1))
	dfa := make(map[int]map[rune]int)
	q := queue.New(1 << uint(n.NumStates))
	q.Push(1 << uint(n.EntryState-1))
	final_states := make([]int, 0)
	for !q.Empty() {
		node_code, _ := q.Pop()
		is_final_node := false
		for node := 1; node < (1 << uint(n.NumStates)); node++ {
			if 1<<uint(node-1) > node_code {
				break
			}
			if ((1 << uint(node-1)) & node_code) == 0 {
				continue
			}
			if is_final[node] {
				is_final_node = true
			}
			if _, ok := dfa[node_code]; !ok {
				dfa[node_code] = make(map[rune]int)
			}
			for character, neighbours := range n.Graph[node] {
				for _, neighbour := range neighbours {
					dfa[node_code][character] = dfa[node_code][character] | (1 << uint(neighbour-1))
					if !visited[dfa[node_code][character]] {
						q.Push(dfa[node_code][character])
						visited[dfa[node_code][character]] = true
					}
				}
			}
		}
		if is_final_node {
			final_states = append(final_states, node_code)
		}
	}

	// encode the newly created DFA into our standard form
	for node, _ := range n.Graph {
		delete(n.Graph, node)
	}
	nr := n.NumStates
	n.NumStates = 0
	n.NumTransitions = 0
	n.EntryState = 1 << uint(n.EntryState-1)

	mapping := make(map[int]int)
	current_node := 0
	for i := 1; i < (1 << uint(nr)); i++ {
		if _, ok := dfa[i]; !ok || len(dfa[i]) == 0 {
			continue
		}

		if _, ok := mapping[i]; !ok {
			current_node++
			mapping[i] = current_node
		}
		n.Graph[mapping[i]] = make(map[rune][]int)
		n.NumTransitions += len(dfa[i])
		for character, node := range dfa[i] {
			if _, ok := mapping[node]; !ok {
				current_node++
				mapping[node] = current_node
			}
			n.Graph[mapping[i]][character] = []int{mapping[node]}
		}
	}
	if _, ok := mapping[n.EntryState]; ok {
		n.EntryState = mapping[n.EntryState]
	} else {
		current_node++
		mapping[n.EntryState] = current_node
		n.EntryState = current_node
	}
	n.NumStates = current_node
	n.FinalStates = make([]int, 0)

	for _, node := range final_states {
		if _, ok := mapping[node]; ok {
			n.FinalStates = append(n.FinalStates, mapping[node])
		}
	}

	return n.DFA
}
