package opt_test

import (
	"slices"
	"testing"

	"github.com/ngicks/option/opt"
)

// cloner is a value type whose Clone performs a deep copy of its slice.
type cloner struct {
	s []int
}

func (c cloner) Clone() cloner {
	return cloner{s: slices.Clone(c.s)}
}

func TestCloneCloner(t *testing.T) {
	orig := opt.Some(cloner{s: []int{1, 2, 3}})
	cloned := opt.CloneCloner(orig)

	if !cloned.IsSome() {
		t.Fatalf("CloneCloner(Some) = %v, want some", cloned)
	}
	if !slices.Equal(cloned.Value().s, []int{1, 2, 3}) {
		t.Errorf("cloned value = %v, want {1,2,3}", cloned.Value().s)
	}

	// Mutating the original's backing array must not affect the clone.
	orig.Value().s[0] = 99
	if cloned.Value().s[0] != 1 {
		t.Errorf("clone aliases original backing array: got %d, want 1", cloned.Value().s[0])
	}

	if got := opt.CloneCloner(opt.None[cloner]()); got.IsSome() {
		t.Errorf("CloneCloner(None) = %v, want none", got)
	}
}
