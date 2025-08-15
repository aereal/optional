package optional

// Unwrap unwraps the contained value from the option.
//
// If the option is none, the false is returned.
func Unwrap[T any](o Option[T]) (T, bool) {
	if o.present {
		return o.v, true
	}
	var t T
	return t, false
}

// IsSome predicates the option is existing value or not.
func IsSome[T any](o Option[T]) bool {
	_, isSome := Unwrap(o)
	return isSome
}

// IsNone predicates the option is none or not.
func IsNone[T any](o Option[T]) bool {
	_, isSome := Unwrap(o)
	return !isSome
}

// Or returns the first existing value or none if all of opts are none.
func Or[T any](opts ...Option[T]) Option[T] {
	for _, o := range opts {
		for v := range o.Iter() {
			return Some(v)
		}
	}
	return None[T]()
}

// Equal returns true if x and y are equal.
//
// The Option's equality is defined by as below:
//
//   - both options are none value, they are equal.
//   - both options are some (existing) value and their underlying value are equal, they are equal.
//   - otherwise, they are NOT equal.
func Equal[T comparable](x, y Option[T]) bool {
	if IsNone(x) && IsNone(y) {
		return true
	}
	if IsSome(x) && IsSome(y) {
		return x.v == y.v
	}
	return false
}
