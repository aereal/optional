package optional

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"iter"
)

// Option represents optional values.
type Option[T any] struct {
	v       T
	present bool
}

var (
	_ json.Marshaler   = Option[any]{}
	_ json.Unmarshaler = (*Option[any])(nil)
	_ driver.Valuer    = Option[any]{}
	_ sql.Scanner      = (*Option[any])(nil)
)

var null = []byte(`null`)

func (o Option[T]) MarshalJSON() ([]byte, error) {
	v, ok := Unwrap(o)
	if !ok {
		return null, nil
	}
	return json.Marshal(v)
}

func (o *Option[T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		*o = None[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}

func (o Option[T]) Value() (driver.Value, error) {
	if !o.present {
		return nil, nil
	}
	return driver.DefaultParameterConverter.ConvertValue(o.v)
}

func (o *Option[T]) Scan(src any) error {
	var n sql.Null[T]
	if err := n.Scan(src); err != nil {
		return err
	}
	*o = FromExistenceCheck(n.V, n.Valid)
	return nil
}

// Iter iterates the option's value.
//
// If the option has no value, nothing iterated.
func (o Option[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.present {
			yield(o.v)
		}
	}
}

// Ptr returns a pointer to the value.
//
// If the option has no value, nil returned.
func (o Option[T]) Ptr() *T {
	if o.present {
		return &o.v
	}
	return nil
}

// FromPtr returns a some value if the pointer refers the existing value.
func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// FromResult returns an existing value if err == nil, otherwise returns a none.
func FromResult[T any](v T, err error) Option[T] {
	if err != nil {
		return None[T]()
	}
	return Some(v)
}

// FromExistenceCheck returns an existing value if exists == true, otherwrise returns a none.
//
// The existence check means (T, bool) tuple values, it is inspired by map type's special index expression.
func FromExistenceCheck[T any](v T, exists bool) Option[T] {
	if exists {
		return Some(v)
	}
	return None[T]()
}

// FromIterator returns an Option with the first value yielded from the iterator.
//
// If the iterator yields no values, returns None.
func FromIterator[T any](i iter.Seq[T]) Option[T] {
	for v := range i {
		return Some(v)
	}
	return None[T]()
}

// Some returns an existing value of type T.
func Some[T any](v T) Option[T] {
	return Option[T]{v: v, present: true}
}

// None returns a none value.
func None[T any]() Option[T] { return Option[T]{present: false} }
