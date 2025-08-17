package optional_test

import (
	"fmt"

	"github.com/aereal/optional"
)

func Example() {
	fmt.Print("Some: ")
	printOption(optional.Some("abc"))
	fmt.Print("None: ")
	printOption(optional.None[string]())
	// Output:
	// Some: "abc" true
	// None: "" false
}

func ExampleFromPtr() {
	s := "abc"
	fmt.Print("presented pointer: ")
	printOption(optional.FromPtr(&s))
	fmt.Print("nil: ")
	printOption(optional.FromPtr((*string)(nil)))
	// Output:
	// presented pointer: "abc" true
	// nil: "" false
}

func ExampleOption_Iter() {
	fmt.Print("Some:")
	for s := range optional.Some("abc").Iter() {
		fmt.Printf(" %q\n", s)
	}
	fmt.Print("None:")
	for s := range optional.None[string]().Iter() {
		fmt.Printf(" %q\n", s)
	}
	// Output:
	// Some: "abc"
	// None:
}

func printOption(o optional.Option[string]) {
	s, ok := optional.Unwrap(o)
	fmt.Printf("%q %v\n", s, ok)
}
