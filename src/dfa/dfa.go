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
	fmt.Printf("\n\n\n\n\nAAA\n\n\n\n")
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
	if len(d.Graph[state]) == 0 {
		return true
	}
	if d.IsFinal(state) {
		return false
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
				match := true
				for character, nodes := range old_graph[i] {
					nodes2, ok := old_graph[j][character]
					if !ok || renames[nodes2[0]] != renames[nodes[0]] {
						match = false
						break
					}
				}
				if match {
					// join states i and j
					obsolete := 0
					used := 0
					if renames[i] > renames[j] {
						obsolete, used = i, j
						renames[i] = renames[j]
						if !to_remove[i] {
							to_remove[i] = true
							nodes_removed++
						}
					} else {
						obsolete, used = j, i
						renames[j] = renames[i]
						if !to_remove[j] {
							to_remove[j] = true
							nodes_removed++
						}
					}
					for color, nodes := range d.Graph[obsolete] {
						d.Graph[used][color][0] = renames[nodes[0]]
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
		for color, _ := range d.Graph[i] {
			for d.Graph[i][color][0] != renames[d.Graph[i][color][0]] {
				d.Graph[i][color][0] = renames[d.Graph[i][color][0]]
			}
		}
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

	// remove/rename the states found
	for node, _ := range d.Graph {
		if to_remove[node] {
			transitions_removed += len(d.Graph[node])
			delete(d.Graph, node)
			continue
		}
		for character, neighbours := range d.Graph[node] {
			if to_remove[neighbours[0]] {
				transitions_removed += len(d.Graph[node][character])
				delete(d.Graph[node], character)
				continue
			}
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
		if mapping[id] != id {
			delete(d.Graph, id)
		}
	}
	d.NumStates -= nodes_removed
	d.NumTransitions -= transitions_removed
}
