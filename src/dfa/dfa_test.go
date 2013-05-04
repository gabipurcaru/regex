package dfa

import (
	"strings"
	"testing"
)

var (
	simple_dfa = "5 4\n" +
		"1 2 l\n" +
		"2 3 o\n" +
		"1 4 a\n" +
		"4 5 l\n" +
		"1\n" +
		"2 3 5"

	complex_dfa = "8 16\n" +
		"1 2 0\n" +
		"1 6 1\n" +
		"2 7 0\n" +
		"2 3 1\n" +
		"3 1 0\n" +
		"3 3 1\n" +
		"4 3 0\n" +
		"4 7 1\n" +
		"5 8 0\n" +
		"5 6 1\n" +
		"6 3 0\n" +
		"6 7 1\n" +
		"7 7 0\n" +
		"7 5 1\n" +
		"8 7 0\n" +
		"8 3 1\n" +
		"1\n" +
		"1 3\n"
)

func BenchmarkDFAProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dfa := NewDFA()
		dfa.Process(strings.NewReader(simple_dfa))
	}
}

func BenchmarkDFACheck(b *testing.B) {
	dfa := NewDFA()
	dfa.Process(strings.NewReader(simple_dfa))
	for i := 0; i < b.N; i++ {
		if !dfa.Check("al") {
			b.Error("Check failed: ac should give true")
		}
	}
}

func BenchmarkDFAMinimize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dfa := NewDFA()
		dfa.Process(strings.NewReader(complex_dfa))
		b.StartTimer()
		dfa.Minimize()
	}
}

func TestDFACheck(t *testing.T) {
	dfa := NewDFA()
	dfa.Process(strings.NewReader(simple_dfa))
	tests := map[string]bool{
		"lo":  true,
		"al":  true,
		"l":   false,
		"lol": false,
		"ll":  false,
		"alo": false,
	}
	for word, res := range tests {
		if dfa.Check(word) != res {
			t.Errorf("Check failed: %s should give %v", word, res)
		}
	}
}

func TestDFAProcess(t *testing.T) {
	dfa := NewDFA()
	dfa.Process(strings.NewReader(simple_dfa))
	if dfa.EntryState != 1 {
		t.Errorf("Incorrect first state: %d", dfa.EntryState)
	}
	if dfa.NumStates != 5 {
		t.Errorf("Incorrect number of states: %d", dfa.NumStates)
	}
	if dfa.NumTransitions != 4 {
		t.Errorf("Incorrect number of transitions: %d", dfa.NumTransitions)
	}
	if len(dfa.FinalStates) != 2 {
		t.Fatalf("Incorrect number of final states: %d", len(dfa.FinalStates))
	}
	if dfa.FinalStates[0] != 3 || dfa.FinalStates[1] != 5 {
		t.Errorf("Final state incorrect: %d", dfa.FinalStates[0])
	}
}

func TestDFAMinimize(t *testing.T) {
	dfa := NewDFA()
	dfa.Process(strings.NewReader(complex_dfa))
	dfa.Minimize()
	tests := map[string]bool{
		"00110":  true,
		"01":     true,
		"011111": true,
		"10":     true,
		"10001":  true,
		"1101":   false,
		"0101":   false,
	}
	if dfa.NumStates != 5 {
		t.Fatalf("Wrong number of minimized states: %d", dfa.NumStates)
	}
	if dfa.NumTransitions != 10 {
		t.Fatalf("Wrong number of minimized transitions: %d", dfa.NumTransitions)
	}
	for test, res := range tests {
		if dfa.Check(test) != res {
			t.Errorf("Wrong answer: %s", test)
		}
	}
}
