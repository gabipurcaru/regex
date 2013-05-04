package nfa

import (
	"strings"
	"testing"
)

var (
	simple_nfa = "4 10\n" +
		"1 2 a\n" +
		"1 2 a\n" +
		"1 4 a\n" +
		"1 3 a\n" +
		"2 3 b\n" +
		"2 4 b\n" +
		"4 3 b\n" +
		"4 2 b\n" +
		"1 2 λ\n" +
		"2 3 λ\n" +
		"1\n" +
		"1 4"
)

func BenchmarkNFAToDFA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		nfa := New()
		nfa.Process(strings.NewReader(simple_nfa))
		b.StartTimer()
		nfa.ToDFA()
	}
}

func TestNFAToDFA(t *testing.T) {
	tests := map[string]bool{
		"ababba":   false,
		"abbbbbb":  true,
		"babababa": false,
		"a":        true,
		"ab":       true,
		"b":        true,
	}
	nfa := New()
	nfa.Process(strings.NewReader(simple_nfa))
	nfa.ToDFA()
	for str, res := range tests {
		if nfa.Check(str) != res {
			t.Fatalf("Wrong answer at: %s", str)
		}
	}
}

func TestConcat(t *testing.T) {
	nfa := New()
	nfa.Process(strings.NewReader(simple_nfa))
	nfa.ToDFA()

	res := Concat(nfa, nfa)
	// res.Print(os.Stdout)

	if res.NumStates != 2*nfa.NumStates {
		t.Errorf("Incorrect number of states")
	}
	if res.NumTransitions != 2*nfa.NumTransitions+len(nfa.FinalStates) {
		t.Errorf("Incorrect number of transitions")
	}
}

func TestEither(t *testing.T) {
	nfa := New()
	nfa.Process(strings.NewReader(simple_nfa))
	nfa.ToDFA()

	res := Either(nfa, nfa)
	// res.Print(os.Stdout)

	if res.NumStates != 2*nfa.NumStates+1 {
		t.Errorf("Incorrect number of states")
	}
	if res.NumTransitions != 2*nfa.NumTransitions+2 {
		t.Errorf("Incorrect number of transitions")
	}
}

func TestStar(t *testing.T) {
	nfa := New()
	nfa.Process(strings.NewReader(simple_nfa))
	nfa.ToDFA()

	Star(nfa)
}
