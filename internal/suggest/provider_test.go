package suggest

import "testing"

func TestFirstMatch(t *testing.T) {
	first := Funcs{
		MatchFn: func(Context) bool { return false },
		SuggestFn: func(Context) []Suggestion {
			return []Suggestion{{Label: "first"}}
		},
	}
	second := Funcs{
		MatchFn: func(Context) bool { return true },
		SuggestFn: func(Context) []Suggestion {
			return []Suggestion{{Label: "second"}}
		},
	}

	got := FirstMatch(Context{}, first, second)
	if len(got) != 1 || got[0].Label != "second" {
		t.Fatalf("got %+v, want second", got)
	}
}

func TestFirstMatchNone(t *testing.T) {
	provider := Funcs{
		MatchFn: func(Context) bool { return false },
		SuggestFn: func(Context) []Suggestion {
			return []Suggestion{{Label: "x"}}
		},
	}
	if got := FirstMatch(Context{}, provider); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}
