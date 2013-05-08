package dfa

import (
	"fmt"
	"io"
	"queue"
)

// DFA represents a Deterministic Finite Automaton
type DFA struct {
	// Number of states in the DFA
	NumStates int

	// Number of transition in the DFA (more useful for reading it)
	NumTransitions int

	// DFA Entry point
	EntryState int

	// DFA Final states
	FinalStates []int

	// The actual DFA is kept as a Graph
	Graph map[int]map[rune][]int
}

// NewDFA is an mpty constructor for a DFA. Returns a null DFA (zero states,
// zero transitions)
func NewDFA() DFA {
	return DFA{0, 0, 0, make([]int, 0), make(map[int]map[rune][]int)}
}

// Process reads a DFA from a Reader. The DFA should look like this:
//
//     [NumStates int] [NumTransitions int]
//     [NumStates] times: [from_state int] [to_state int] [character rune]
//     [EntryState int] [num_FinalStates int]
//     [num_FinalStates] times: [state int]
func (d *DFA) Process(r io.Reader) {
	var a, b int
	var c rune
	var NumStates, state int

	fmt.Fscanf(r, "%d %d\n", &d.NumStates, &d.NumTransitions)
	for i := 0; i < d.NumTransitions; i++ {
		_, err := fmt.Fscanf(r, "%d %d %c\n", &a, &b, &c)
		if err != nil {
			panic(err)
		}
		_, ok := d.Graph[a]
		if !ok {
			d.Graph[a] = make(map[rune][]int)
		}
		if _, ok := d.Graph[a][c]; !ok {
			d.Graph[a][c] = []int{b}
		} else {
			d.Graph[a][c] = append(d.Graph[a][c], b)
		}
	}

	fmt.Fscanf(r, "%d %d\n", &d.EntryState, &NumStates)
	for i := 0; i < NumStates; i++ {
		fmt.Fscanf(r, "%d\n", &state)
		d.FinalStates = append(d.FinalStates, state)
	}
}

// Checks whether the DFA accepts a given word
func (d *DFA) Check(word string) bool {
	state := d.EntryState
	for _, char := range word {
		nodes, ok := d.Graph[state][char]
		if !ok {
			return false
		}
		state = nodes[0]
	}
	return d.IsFinal(state)
}

// Prints the DFA, using the same format it uses to read it
func (d *DFA) Print(w io.Writer) {
	fmt.Fprintf(w, "%d %d\n", d.NumStates, d.NumTransitions)
	for i := 1; i <= d.NumStates; i++ {
		for character, nodes := range d.Graph[i] {
			for _, neighbour := range nodes {
				fmt.Fprintf(w, "%d %d %c\n", i, neighbour, character)
			}
		}
	}
	fmt.Fprintf(w, "%d\n", d.EntryState)
	fmt.Fprintf(w, "%d ", len(d.FinalStates))
	for _, state := range d.FinalStates {
		fmt.Fprintf(w, "%d ", state)
	}
}

// Returns true if the state is a final DFA state
func (d *DFA) IsFinal(state int) bool {
	for _, node := range d.FinalStates {
		if node == state {
			return true
		}
	}
	return false
}

// Returns true if there isn't a word that takes this state into a final ones
func (d *DFA) IsReverseUnreachable(state int) bool {
	if d.IsFinal(state) {
		return false
	}
	if len(d.Graph[state]) == 0 {
		return true
	}
	q := queue.New(d.NumStates)
	viz := make([]bool, d.NumStates+1)
	q.Push(state)
	for !q.Empty() {
		node, _ := q.Pop()
		for _, neighbours := range d.Graph[node] {
			if viz[neighbours[0]] {
				continue
			}
			if d.IsFinal(neighbours[0]) {
				return false
			}
			viz[neighbours[0]] = true
			q.Push(neighbours[0])
		}
	}
	return true
}

func nodes_match(graph map[int]map[rune][]int, is_final []bool, node1, node2 int) bool {
	// perform a simultaneous BFS from the two nodes; if there is a string
	// that matches from node1 but not from node2 or the other way around,
	// the two nodes don't match
	// q := queue.New(1)
	if len(graph[node1]) != len(graph[node2]) {
		return false
	}
	for character, neighbours := range graph[node1] {
		if neighbours2, ok := graph[node2][character]; !ok || is_final[neighbours[0]] != is_final[neighbours2[0]] {
			return false
		}
	}
	return true
}

// Minimize simplifies the DFA, by removing unnecessary states and merging
// states that can be merged
func (d *DFA) Minimize() {
	// find unreachable and reverse-unreachable states
	q := queue.New(d.NumStates)
	to_remove := make([]bool, d.NumStates+1)
	viz := make(map[int]bool)
	nodes_removed := 0
	transitions_removed := 0
	q.Push(d.EntryState)
	viz[d.EntryState] = true
	for !q.Empty() {
		node, _ := q.Pop()
		for _, neighbours := range d.Graph[node] {
			ok, _ := viz[neighbours[0]]
			if !ok {
				viz[neighbours[0]] = true
				q.Push(neighbours[0])
			}
		}
	}
	for i := 1; i <= d.NumStates; i++ {
		ok, _ := viz[i]
		if !ok || d.IsReverseUnreachable(i) {
			to_remove[i] = true
			nodes_removed++
		}
	}

	// Moore's Algorithm
	// See http://en.wikipedia.org/wiki/DFA_minimization#Moore.27s_algorithm
	is_final := make([]bool, d.NumStates+1)
	renames := make(map[int]int)
	for i := 1; i <= d.NumStates; i++ {
		renames[i] = i
		if d.IsFinal(i) {
			is_final[i] = true
		}
	}
	found_match := true
	for found_match {
		found_match = false
		old_graph := d.Graph
		for i := 1; i <= d.NumStates; i++ {
			if to_remove[renames[i]] {
				continue
			}
			for j := 1; j <= d.NumStates; j++ {
				if to_remove[renames[i]] || renames[i] >= renames[j] {
					continue
				}
				if is_final[renames[i]] && !is_final[renames[j]] ||
					is_final[renames[j]] && !is_final[renames[i]] {
					continue
				}
				match := nodes_match(old_graph, is_final, i, j)
				if match {
					// join states i and j
					obsolete := 0
					used := 0
					if renames[i] > renames[j] {
						obsolete, used = renames[i], renames[j]
						renames[i] = renames[j]
					} else {
						obsolete, used = renames[j], renames[i]
						renames[j] = renames[i]
					}
					if !to_remove[obsolete] {
						to_remove[obsolete] = true
						nodes_removed++
					}
					for character, nodes := range d.Graph[obsolete] {
						if _, ok := d.Graph[used]; !ok {
							d.Graph[used] = make(map[rune][]int)
						}
						if _, ok := d.Graph[used][character]; !ok {
							d.Graph[used][character] = make([]int, 1)
						} else {
							transitions_removed++
						}
						d.Graph[used][character][0] = renames[nodes[0]]
					}
					renames[obsolete] = renames[used]
					found_match = true
				}
			}
		}
	}

	// Apply remaining renames
	// TODO: improve the main algorithm so that it doesn't need this
	for i := 1; i <= d.NumStates; i++ {
		for character, _ := range d.Graph[i] {
			for d.Graph[i][character][0] != renames[d.Graph[i][character][0]] {
				d.Graph[i][character][0] = renames[d.Graph[i][character][0]]
			}
		}
	}
	for d.EntryState != renames[d.EntryState] {
		d.EntryState = renames[d.EntryState]
	}

	// rename states so that there are no gaps between them
	mapping := make(map[int]int)
	next_available := 0
	for node := 1; node <= d.NumStates; node++ {
		if !to_remove[node] {
			next_available++
			mapping[node] = next_available
		}
	}
	d.EntryState = mapping[d.EntryState]
	final_states := make([]int, 0, len(d.FinalStates))
	for id, _ := range d.FinalStates {
		if _, ok := mapping[d.FinalStates[id]]; ok {
			final_states = append(final_states, mapping[d.FinalStates[id]])
		}
	}
	d.FinalStates = final_states

	// rename the states found
	for node, _ := range d.Graph {
		for character, _ := range d.Graph[node] {
			d.Graph[node][character][0] = mapping[d.Graph[node][character][0]]
		}
	}
	for id := 1; id <= d.NumStates; id++ {
		_, ok := mapping[id]
		if !ok {
			delete(d.Graph, id)
			continue
		}
		d.Graph[mapping[id]] = d.Graph[id]
	}
	d.NumStates -= nodes_removed
	d.NumTransitions -= transitions_removed
}
