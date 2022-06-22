# Code Style Guide

* [Basic Rules](#basic-rules)
* [Linter](#linter)
* [Unit Tests](#unit-tests)
* [Additional Rules](#additional-rules)

## Basic Rules

As a basis, we use the style of Uber: https://github.com/uber-go/guide/blob/master/style.md.

Some clarifications, changes and additions to the Uber guide:

- **“Error Wrapping”** :warning: Use `%w` instead of `%v` to wrap an error.
- **“Use go.uber.org/atomic”** :no_entry:  
  We will use `go.uber.org/atomic` only if we have to work with different types of atomics.
- **“Avoid Mutable Globals”** :no_entry:  
  There is another way to solve this problem which is used inside standard packages:

  ```go
  func Sign(msg string) string {
    return defaultSigner.Sign(msg)
  }

  var defaultSigner = &signer{
    now: time.Now,
  }

  type signer struct {
    now func() time.Time
    
    // mx sync.Mutex if we need thread safety
  }

  func (s *signer) Sign(msg string) string {
    return signWithTime(msg, s.now())
  }
  ```

- **“Performance”** :warning:  
  This section should be treated with healthy skepticism and confirm the hypotheses with benchmarks.
- **“Function Grouping and Ordering”** :warning:  
  In general, it is a good example, but if a private function is used only within one function, then
  it is better to place its definition immediately after it. It makes the reading easier, and there
  is no need to run through the file from one place to another.

  ```go
  func Foo() {
    // ...
  
    bar()
  
    // ...
  }
  
  func bar() {
    // ...
  }
  ```

- **“Prefix Unexported Globals with _`”** :no_entry:  
  Ignore this section entirely as it is weird.

## Linter

To check the source code, [golangci-lint](https://golangci-lint.run/) linter is used.

### Installing the **golangci-lint** linter

1. Install the linter by following the instructions:
   https://golangci-lint.run/usage/install/#local-installation
2. Run the linter at the project root with the command `golangci-lint run`.

The linter configuration is located in the file `.golangci.yml` in the project root.

More details about the file format can be found here: https://golangci-lint.run/usage/configuration/  
See available linters and their settings here: https://golangci-lint.run/usage/linters/

To ease the launch, a target `lint` has been added to the `Makefile`.

## Unit Tests

Unit test format:

```go
func TestFoo(t *testing.T) {
  t.Run("positive", func(t *testing.T) {
    t.Run("case #1", func(t *testing.T) { ... })
    t.Run("case #2", func(t *testing.T) { ... })
    ...
  })
  t.Run("negative", func(t *testing.T) {
    t.Run("case #1", func(t *testing.T) { ... })
    t.Run("case #2", func(t *testing.T) { ... })
    ...
  })
}
```

Thus, all tests of one function/method will be collected in one place.

It is also recommended:

- To use https://github.com/stretchr/testify in unit tests where you will find many ready-made
  primitives for asserts.
- To use table tests https://github.com/golang/go/wiki/TableDrivenTests.

## Additional Rules

- Do not use `.` imports. The package name is important in understanding what’s going on.
- A good approach, if possible, is to describe the entire public interface in one file, the name of
  which is the same as the package name.
- Use the standard approach to document public symbols, where the symbol name is specified first in
  the comment followed by its description after a space.

  ```go
  // Foo do some staff.
  func Foo() {
  ```

- The package description should come before the keyword `package` and be in the format `Package
  <package name> <package description>`. For example:
  
  ```go
  // Package foo is designed for ...
  package foo
  ```

- Validator functions should return an error, not panic. This will simplify the code. In general, it
  is worth getting rid of `panic` and `recover` in the code.

- In each project, add a `Makefile` with multiple targets:
  - **build** to build if it's an application;
  - **fmt** to format code;
  - **lint** to run linters;
  - **test** to run tests.
