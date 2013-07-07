package regex

import (
	"os"
	"testing"
)

func TestRegexCheck(t *testing.T) {
	type Test struct {
		Re      string
		Match   string
		Matches bool
	}
	tests := []Test{
		Test{"(a|b)*blabla", "abababblabla", true},
		Test{"(a|b)*blabla", "abababblab", false},
		Test{"(a|bb)*", "abbaaaabbbb", true},
		Test{"(a|bb)*", "", true},
		Test{"(a|bb)*", "aaaaaa", true},
		Test{"(a|bb)*", "abbab", false},
		Test{"a|b|c", "b", true},
		Test{"((a|b*)c)*", "accccccc", true},
		Test{"((a|b*)c)*", "aaccccc", false},
		Test{"((a|b*)c)*", "abbbbbbccccbbccbcbcbcabc", false},
	}
	for _, test := range tests {
		nfa := RegexToNFA(test.Re)
		dfa := nfa.ToDFA()
		dfa.Minimize()
		if dfa.Check(test.Match) != test.Matches {
			t.Errorf("Regex fails: %s with %s", test.Re, test.Match)
			dfa.Print(os.Stderr)
		}
	}
}
