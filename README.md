# argparse
## Project Status
[ ![Codeship Status for clagraff/argparse](https://codeship.com/projects/68eb7800-af6b-0133-1b97-3e80188314d9/status?branch=master)](https://codeship.com/projects/132507)
[![GoDoc](https://godoc.org/github.com/clagraff/argparse?status.svg)](https://godoc.org/github.com/clagraff/argparse)

## Description
A Golang flag & argument parser for program arguments. The goal of argparse (program-arguments) is to emulate the usability and functionallity of Python's [argparse](https://docs.python.org/dev/howto/argparse.html#the-basics) package, where setting up a parser and arguments is both easy and trivial.

# Install
Installing is simple, as with most other Golang packages:

```bash
$ go get github.com/clagraff/argparse
```

Boom! All set! Feel free to read on for examples on getting started and using this package.

# The Basics
## Create a parser
Here we have a basic program which greets a user by a provided name, optionally in uppercase. We need create a parser, include our option and flag, and then parse our programs arguments. It could look something like:

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/clagraff/argparse"
)

func main() {
	p := argparse.NewParser("Output a friendly greeting").Version("1.3.0a")
	p.AddHelp().AddVersion() // Enable help and version flags

	upperFlag := argparse.NewFlag("u", "upper", "Use uppercase text").Default("false")
	nameOption := argparse.NewArg("n name", "name", "Name of person to greet").Default("John").Required()

	p.AddOptions(dry_run, max)

	// Parse all available program arguments (except for the program path).
	if ns, leftovers, err := p.Parse(os.Args[1:]...); err != nil {
		switch err.(type) {
		case argparse.ShowHelpErr, argparse.ShowVersionErr:
		    // For either ShowHelpErr or ShowVersionErr, the parser has already
		    // displayed the necessary text to the user. So we end the program 
		    // by returning.
			return
		default:
			fmt.Println(err, "\n")
			p.ShowHelp()
		}
	} else {
		name := ns.String("name")
		upper := ns.String("upper") == "true"

		if upper == true {
			name = strings.ToUpper(name)
		}

		fmt.Printf("Hello, %s!\n", name)
		if len(leftovers) > 0 {
			fmt.Println("\nUnused args:", leftovers)
		}
	}
}
```

You could then run it and receive the following output:

```bash
> go run main.go Luke
Hello, Luke!

> go run main.go Vader -u
Hello, VADER!

> go run main.go
n, name: too few arguments

usage: main [-h] [-v] [-u] [n NAME]

Output a friendly greeting

positional arguments:
  [n NAME]       Name of person to greet

optional arguments:
  -h, --help     Show program help
  -v, --version  Show program version
  -u, --upper    Use uppercase text
```

## Arguments
Arguments are command-line values passed to the program when its execution starts. When these
values are expected by the program, we use a convention of classifying these arguments
into two types: Flags and Options.

### Types
#### Flags
Flags represent non-positional, boolean arguments. These arguments are `"false"` 
by default, will utilize the `argparse.StoreTrue` action. They do not consume 
any additional arguments other than themselves. You can create a new flag using:

```go
// Create a short and long flag for "use-default", -d --default, with a description.
use_default := argparse.NewFlag("d default", "use-default", "Enable the default mode")
```

`argparse.Option` is the struct used for creating parseable options. 

#### Options
Options are arguments which store/represent one or more values. While options 
can represent a variety of input types, they are always serialized to
a `string`. Options may be required or not, may be positional or not, can have
a variable number of parameters, and can operate in a variety of manners.

Here is an example of an option `foo`, which has a default value of `bar`, 
is required, and expects 1 additional argument which will be stored as `foo`'s value
```go
// Create a required Option -foo, with a default value of "bar", that expects and stores 1 argument.
f := argparse.NewOption("-f --foo", "foo", "A foo option").Default("bar").Nargs("1").Required().Action(argparse.Store)
```

### Methods
Options can be configured in a variety of ways. Therefore, method-chaining is
heavily used to quickly create and setup an option. Consider the following example:

```go
// Create a required Option -foo, with a default value of "bar", that expects and stores 1 argument.
f := argparse.NewOption("-f --foo", "foo", "A foo option").Default("bar").Nargs("1").Required().Action(argparse.Store)
```

Options can have the following attributes:
* Is required or not required
* Is positional or not positional
* Has a default value
* Has a constant value
* Expects a specified number of arguments (or no arguments)
* Is identified by one or more public qualifiers (e.g.: `-f` or `--foo`)
* Can require arguments to match specified choices

#### Nargs
Nargs, a shortening of "numer of arguments", represents the number of arguments a flag expects after its presence in a programs complete list of parameters. This could be an actual number, such as `0` or `5`, or it could be any of the following characters: `*+?`. 

The `*` character represents "any and all arguments" following the flag.

The `+` character represents "one or more arguments" following the flag.

The `?` character represents "no arguments or only one argument" following the flag.

The `r` or `R` characters represent "all remaining arguments" that were not consumed during parsing. These narg choices do not consume the parse arguments they are applicable to.

#### Actions
A flag's `action` defines what should occur when a flag is parsed. All flags must have an action. By default, a flag will store `true` in the parser when present, and `false` when not. The following are the currently available actions:

* __argparse.StoreTrue__ will store `true` in the parser when the flag is present.
* __argparse.StoreFalse__ will store `false` in the parser when the flag is present.
* __argparse.StoreConst__ will store the flag's constant value into the parser when the flag is present.
* __argparse.Store__ will store the appropriate number of arguments into the parser when the flag & arguments are present.
* __argparse.AppendConst__ will append the flag's constant to the flag's slice within the parser.
* __argparse.Append__ will append the appropriate number of arguments into the flag's slice within the parser.
* __argparse.ShowHelp__ will print the parser's generate help text to `stdout`.
