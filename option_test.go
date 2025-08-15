package optional_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/aereal/optional"
)

var (
	n       = 10
	some    = optional.FromPtr(&n)
	none    = optional.FromPtr[int](nil)
	errOops = errors.New("oops")
)

func TestFromPtr(t *testing.T) {
	t.Parallel()

	if optional.IsNone(some) || !optional.IsSome(some) {
		t.Error("expected FromPtr(&n) returns a some")
	}

	if !optional.IsNone(none) || optional.IsSome(none) {
		t.Error("expected FromPtr(nil) returns a none")
	}
}

func TestFromResult(t *testing.T) {
	t.Parallel()
	t.Run("err == nil", func(t *testing.T) {
		t.Parallel()
		got := optional.FromResult(succeeds())
		assertsSome(t, got, 10)
	})
	t.Run("err != nil", func(t *testing.T) {
		t.Parallel()
		got := optional.FromResult(fails())
		assertsNone(t, got)
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()
	t.Run("some", func(t *testing.T) {
		t.Parallel()
		v, ok := optional.Unwrap(some)
		if !ok {
			t.Fatal("unwrap failed")
		}
		if v != 10 {
			t.Errorf("unexpected unwrapped value: %d", v)
		}
	})
	t.Run("none", func(t *testing.T) {
		t.Parallel()
		_, ok := optional.Unwrap(none)
		if ok {
			t.Error("unexpectedly unwrap succeeds")
		}
	})
}

func TestOr(t *testing.T) {
	t.Parallel()
	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		got := optional.Or[int]()
		assertsNone(t, got)
	})
	t.Run("none only", func(t *testing.T) {
		t.Parallel()
		got := optional.Or(none, none)
		assertsNone(t, got)
	})
	t.Run("[none, some1, some2]", func(t *testing.T) {
		t.Parallel()
		got := optional.Or(none, some, optional.Some(20))
		assertsSome(t, got, 10)
	})
}

func TestEqual(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		lhs, rhs optional.Option[int]
		want     bool
	}{
		{
			lhs:  optional.Some(123),
			rhs:  optional.Some(123),
			want: true,
		},
		{
			lhs:  optional.Some(123),
			rhs:  optional.Some(456),
			want: false,
		},
		{
			lhs:  optional.Some(123),
			rhs:  none,
			want: false,
		},
		{
			lhs:  none,
			rhs:  none,
			want: true,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v vs %v", tc.lhs, tc.rhs), func(t *testing.T) {
			t.Parallel()
			got := optional.Equal(tc.lhs, tc.rhs)
			if got != tc.want {
				t.Errorf("want=%v; got=%v", tc.want, got)
			}
		})
	}
}

func TestOption_Value(t *testing.T) {
	t.Parallel()
	gotSome := slices.Collect(some.Iter())
	if got := gotSome[0]; got != 10 {
		t.Errorf("some() value returns unexpected value: %d", got)
	}

	gotNone := slices.Collect(none.Iter())
	if len(gotNone) > 0 {
		t.Errorf("none() value returns unexpected value: %#v", gotNone)
	}
}

func TestOption_Ptr(t *testing.T) {
	t.Parallel()
	p1 := some.Ptr()
	if p1 == nil || *p1 != 10 {
		t.Errorf("some.Ptr() returns unexpected value: %#v", p1)
	}

	p2 := none.Ptr()
	if p2 != nil {
		t.Errorf("none.Ptr() returns unexpected value: %#v", p2)
	}
}

func TestOption_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("some[int]", func(t *testing.T) {
		t.Parallel()
		got, err := json.Marshal(some)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "10" {
			t.Errorf("unexpected JSON marshaling: %q", string(got))
		}
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()
		got, err := json.Marshal(none)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "null" {
			t.Errorf("unexpected JSON marshaling: %q", string(got))
		}
	})

	t.Run("some[func]", func(t *testing.T) {
		t.Parallel()
		_, err := json.Marshal(optional.Some(func() {}))
		if err == nil {
			t.Error("expected some error but got nil")
		}
	})
}

func TestOption_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		wantErr error
		input   []byte
		want    optional.Option[int]
	}{
		{
			input:   []byte(`123`),
			want:    optional.Some(123),
			wantErr: nil,
		},
		{
			input:   []byte(`null`),
			want:    optional.None[int](),
			wantErr: nil,
		},
		{
			input:   []byte(`"abc"`),
			want:    optional.None[int](),
			wantErr: literalError("json: cannot unmarshal string into Go value of type int"),
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input: %s", string(tc.input)), func(t *testing.T) {
			t.Parallel()
			var got optional.Option[int]
			gotErr := json.Unmarshal(tc.input, &got)
			if !errors.Is(tc.wantErr, gotErr) {
				t.Errorf("error:\n\twant: %s (%T)\n\t got: %s (%T)", tc.wantErr, tc.wantErr, gotErr, gotErr)
			}
			if !optional.Equal(tc.want, got) {
				t.Errorf("Option:\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

func assertsSome[T comparable](t *testing.T, o optional.Option[T], want T) {
	t.Helper()

	got, ok := optional.Unwrap(o)
	if !ok {
		t.Fatal("expected some value but got none")
	}
	if got != want {
		t.Errorf("value mismatch:\n\twant: %#v\n\t got: %#v", want, got)
	}
}

func assertsNone[T any](t *testing.T, o optional.Option[T]) {
	t.Helper()

	if !optional.IsNone(o) {
		t.Fatal("unexpectedly unwrap suceeds")
	}
}

func succeeds() (int, error) { return 10, nil }

func fails() (int, error) { return 0, errOops }

type literalError string

var _ error = (literalError)("")

func (e literalError) Error() string { return string(e) }

func (e literalError) Is(other error) bool {
	return e.Error() == other.Error()
}
