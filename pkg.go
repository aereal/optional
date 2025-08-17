// Package optional provides a generic optional type for Go, inspired by
// [Scala's Option], and [Haskell's Maybe] types.
//
// It offers a type-safe way to handle values that may or may not exist,
// eliminating the need for nil pointer checks and providing clearer semantics than
// (T, error) patterns or bare nil values.
//
// The core type [Option] encapsulates an optional value of type T.
// Unlike bare pointers or any values, Option[T] makes the absence of a value explicit
// and type-safe, preventing runtime panics from nil dereferences.
//
// # Safety and Semantics
//
// [Option] provides safer and more explicit semantics compared to traditional Go patterns:
//   - No nil pointer dereferences: values are accessed through safe methods
//   - Clear intention: the type system enforces handling of absent values
//
// # Integration
//
// [Option] integrates seamlessly with Go's standard library:
//   - JSON marshaling: implements [json.Marshaler] and [json.Unmarshaler]; Some(value) → value, None() → null
//   - SQL databases: implements [driver.Valuer] and [sql.Scanner] interfaces
//   - Iteration: works with Go 1.23+ [iterators] via [Option.Iter] method
//
// # Memory Efficiency
//
// [Option] is designed with memory efficiency in mind:
//   - Zero allocation for None values after initial creation
//   - Single allocation for Some values
//   - Compact representation: stores value and presence flag inline
//
// # Extensibility
//
// The iterator-based design enables composition using standard iterator utilities.
//
// This allows complex optional value manipulations using familiar Go iteration patterns without
// requiring [Option]-specific utility functions.
//
// [Scala's Option]: https://www.scala-lang.org/api/current/scala/Option.html
// [Haskell's Maybe]: https://hackage.haskell.org/package/base-4.21.0.0/docs/Data-Maybe.html
// [iterators]: https://pkg.go.dev/iter#hdr-Iterators
package optional

import (
	_ "database/sql"
	_ "database/sql/driver"
	_ "encoding/json"
)
