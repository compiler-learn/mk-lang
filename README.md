# monkey-lang

![](https://github.com/prologic/monkey-lang/workflows/Coverage/badge.svg)
![](https://github.com/prologic/monkey-lang/workflows/Docker/badge.svg)
![](https://github.com/prologic/monkey-lang/workflows/Go/badge.svg)
![](https://github.com/prologic/monkey-lang/workflows/ReviewDog/badge.svg)

[![CodeCov](https://codecov.io/gh/prologic/monkey-lang/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/monkey-lang)
[![Go Report Card](https://goreportcard.com/badge/prologic/monkey-lang)](https://goreportcard.com/report/prologic/monkey-lang)
[![GoDoc](https://godoc.org/github.com/prologic/monkey-lang?status.svg)](https://godoc.org/github.com/prologic/monkey-lang) 
[![Sourcegraph](https://sourcegraph.com/github.com/prologic/monkey-lang/-/badge.svg)](https://sourcegraph.com/github.com/prologic/monkey-lang?badge)

Monkey programming language interpreter designed in [_Writing An Interpreter In Go_](https://interpreterbook.com/).
A step-by-step walk-through where each commit is a fully working part.
Read the book and follow along with the commit history.

## Table of Contents

* [monkey\-lang](#monkey-lang)
  * [Status](#status)
  * [Read and Follow](#read-and-follow)
  * [Quickstart](#quickstart)
  * [Development](#development)
  * [Monkey Language](#monkey-language)
    * [Programs](#programs)
    * [Types](#types)
    * [Variable Bindings](#variable-bindings)
    * [Artithmetic Expressions](#artithmetic-expressions)
    * [Conditional Expressions](#conditional-expressions)
    * [While Loops](#while-loops)
    * [Functions and Closures](#functions-and-closures)
    * [Recursive Functions](#recursive-functions)
    * [Strings](#strings)
    * [Arrays](#arrays)
    * [Hashes](#hashes)
    * [Assignment Expressions](#assignment-expressions)
    * [Binary and unary operators](#binary-and-unary-operators)
    * [Builtin functions](#builtin-functions)
    * [Objects](#objects)
    * [Modules](#modules)
  * [License](#license)

## Status

> Still working on a self-hosted Monkey lang (*Monkey written in Monkey*).

## Read and Follow

> Read the books and follow along with the following commit history.
(*This also happens to be the elapsed days I took to read both books!*)

See: [Reading Guide](./ReadingGuide.md)

> Please note that whilst reading the awesome books I slightly modified this
> version of Monkey-lang in some places. For example I opted to have a single
> `RETURN` Opcode.

## Quick start

```#!sh
$ go get github.com/prologic/monkey-lang
$ monkey-lang
```

## Development

To build run `make`.

```#!sh
$ git clone https://github.com/prologic/monkey-lang
$ monkey-lang
$ make
This is the Monkey programming language!
Feel free to type in commands
>> 
```

To run the tests run `make test`

You can also execute program files by invoking `monkey-lang <filename>`
There are also some command-line options:

```#!sh
$ ./monkey-lang -h
Usage: monkey-lang [options] [<filename>]
  -c	compile input to bytecode
  -d	enable debug mode
  -e string
    	engine to use (eval or vm) (default "vm")
  -i	enable interactive mode
  -v	display version information
```

## Monkey Language

> See also: [examples](./examples)

### Programs

A Monkey program is simply zero or more statements. Statements don't actually
have to be separated by newlines, only by white space. The following is a valid
program (*but you'd probably use newlines in the`if` block in real life*):

```
s := "world"
print("Hello, " + s)
if (s != "") { t := "The end" print(t) }
// Hello, world
// The end
```

Between tokens, white space and comments
(*lines starting with `//` or `#` through to the end of a line*)
are ignored.

### Types

Monkey has the following data types: `null`, `bool`, `int`, `str`, `array`,
`hash`, and `fn`. The `int` type is a signed 64-bit integer, strings are
immutable arrays of bytes, arrays are grow-able arrays
(*use the `append()` builtin*), and hashes are unordered hash maps.
Trailing commas are **NOT** allowed after the last element in an array or hash:

Type      | Syntax                                    | Comments
--------- | ----------------------------------------- | -----------------------
null      | `null`                                    |
bool      | `true false`                              |
int       | `0 42 1234 -5`                            | `-5` is actually `5` with unary `-`
str       | `"" "foo" "\"quotes\" and a\nline break"` | Escapes: `\" \\ \t \r \n \t \xXX`
array     | `[] [1, 2] [1, 2, 3]`                     |
hash      | `{} {"a": 1} {"a": 1, "b": 2}`            |

### Variable Bindings

```#!sh
>> a := 10
```

### Arithmetic Expressions

```#!sh
>> a := 10
>> b := a * 2
>> (a + b) / 2 - 3
12
```

### Conditional Expressions

Monkey supports `if` and `else`:

```sh
>> a := 10
>> b := a * 2
>> c := if (b > a) { 99 } else { 100 }
>> c
99
```

Monkey also supports `else if`:

```#!sh
>> test := fn(n) { if (n % 15 == 0) { return "FizzBuzz" } else if (n % 5 == 0) { return "Buzz" } else if (n % 3 == 0) { return "Fizz" } else { return str(n) } }
>> test(1)
"1"
>> test(3)
"Fizz"
>> test(5)
"Buzz"
>> test(15)
"FizzBuzz"
```

### While Loops

Monkey supports only one looping construct, the `while` loop:

```#!sh
i := 3
while (i > 0) {
    print(i)
    i = i - 1
}
// 3
// 2
// 1
```

Monkey does not have `break` or `continue`, but you can `return <value>` as
one way of breaking out of a loop early inside a function.

### Functions and Closures

You can define named or anonymous functions, including functions inside
functions that reference outer variables (*closures*).

```sh
>> multiply := fn(x, y) { x * y }
>> multiply(50 / 2, 1 * 2)
50
>> fn(x) { x + 10 }(10)
20
>> newAdder := fn(x) { fn(y) { x + y } }
>> addTwo := newAdder(2)
>> addTwo(3)
5
>> sub := fn(a, b) { a - b }
>> applyFunc := fn(a, b, func) { func(a, b) }
>> applyFunc(10, 2, sub)
8
```

**NOTE:** You cannot have a "bare return" -- it requires a return value.
          So if you don't want to return anything
          (*functions always return at least `null` anyway*),
          just say `return null`.

### Recursive Functions

Monkey also supports recursive functions including recursive functions defined
in the scope of another function (*self-recursion*).

```#!sh
>> wrapper := fn() { inner := fn(x) { if (x == 0) { return 2 } else { return inner(x - 1) } } return inner(1) }
>> wrapper()
2
```

Monkey also does tail
call optimization and turns recursive tail-calls into iteration.

```#!sh
>> fib := fn(n, a, b) { if (n == 0) { return a } if (n == 1) { return b } return fib(n - 1, b, a + b) }
>> fib(35, 0, 1)
9227465
```

### Strings

```sh
>> makeGreeter := fn(greeting) { fn(name) { greeting + " " + name + "!" } }
>> hello := makeGreeter("Hello")
>> hello("skatsuta")
Hello skatsuta!
```

### Arrays

```sh
>> myArray := ["Thorsten", "Ball", 28, fn(x) { x * x }]
>> myArray[0]
Thorsten
>> myArray[4 - 2]
28
>> myArray[3](2)
4
```

### Hashes

```sh
>> myHash := {"name": "Jimmy", "age": 72, true: "yes, a boolean", 99: "correct, an integer"}
>> myHash["name"]
Jimmy
>> myHash["age"]
72
>> myHash[true]
yes, a boolean
>> myHash[99]
correct, an integer
```

### Assignment Expressions

Assignment can assign to a name, an array element by index, or a hash value by key.
When assigning to a name (variable), it always assigns to the scope the variable was defined .

To help with object-oriented programming, `obj.foo = bar` is syntactic sugar for `obj["foo"] = bar`. They're exactly equivalent.

```
i := 1
func mutate() {
    i = 2
    print(i)
}
print(i)
mutate()
print(i)
// 1
// 2
// 2

map = {"a": 1}
func mutate() {
    map.a = 2
    print(map.a)
}
print(map.a)
mutate()
print(map.a)
// 1
// 2
// 2

lst := [0, 1, 2]
lst[1] = "one"
print(lst)
// [0, "one", 2]

map = {"a": 1, "b": 2}
map["a"] = 3
map.c = 4
print(map)
// {"a": 3, "b": 2, "c": 4}
```

### Binary and unary operators

Monkey supports pretty standard binary and unary operators.
Here they are with their precedence, from highest to lowest
(*operators of the same precedence evaluate left to right*):

Operators      | Description
-------------- | -----------
`[] obj.keu`   | Subscript
`-`            | Unary minus
`* / %`        | Multiplication, Division, Modulo
`+ -`          | Addition, Subtraction
`< <= > >= in` | Comparison
`== !=`        | Equality
`<< >>`        | Bit Shift
`~`            | Bitwise not
`&`            | Bitwise and
<code>&#124;</code>       | Bitwise or
<code>&#124;&#124;</code> | Logical or (short-circuit)
`&&`           | Logical and (short-circuit)
`!`            | Logical not

Several of the operators are overloaded. Here are the types they can operate on:

Operator   | Types           | Action
---------- | --------------- | ------
`[]`       | `str[int]`      | fetch nth byte of str (0-based)
`[]`       | `array[int]`    | fetch nth element of array (0-based)
`[]`       | `hash[str]`     | fetch hash value by key str
`-`        | `int`           | negate int
`*`        | `int * int`     | multiply ints
`*`        | `str * int`     | repeat str n times
`*`        | `int * str`     | repeat str n times
`*`        | `array * int`   | repeat array n times, give new array
`*`        | `int * array`   | repeat array n times, give new array
`/`        | `int / int`     | divide ints, truncated
`%`        | `int % int`     | divide ints, give remainder
`+`        | `int + int`     | add ints
`+`        | `str + str`     | concatenate strs, give new string
`+`        | `array + array` | concatenate arrays, give new array
`+`        | `hash + hash`   | merge hashes into new hash, keys in right hash win
`-`        | `int - int`     | subtract ints
`<`        | `int < int`     | true iff left < right
`<`        | `str < str`     | true iff left < right (lexicographical)
`<`        | `array < array` | true iff left < right (lexicographical, recursive)
`<= > >=`  | same as `<`     | similar to `<`
`<<`       | `int << int`    | Shift left by n bits
`>>`       | `int >> int`    | Shift right by n bits
`==`       | `any == any`    | deep equality (always false if different type)
`!=`       | `any != any`    | same as `not ==`
<code>&#124;</code> | <code>int &#124; int</code> | Bitwise or
`&`        | `int & int`     | Bitwise and
`~`        | `~int`          | Bitwise not (1's complement)
<code>&#124;&#124;</code> | <code>bool &#124;&#124; bool</code> | true iff either true, right not evaluated if left true
`&&`       | `bool && bool`  | true iff both true, right not evaluated if left false
`!`        | `!bool`         | inverse of bool

### Builtin functions

- `len(iterable)`
  Returns the length of the iterable (`str`, `array` or `hash`).
- `input([prompt])`
  Reads a line from standard input optionally printing `prompt`.
- `print(value...)`
  Prints the `value`(s) to standard output followed by a newline.
- `first(array)`
  Returns the first element of the `array`.
- `last(array)`
  Returns the last element of the `array`.
- `rest(array)`
  Returns a new array with the first element of `array` removed.
- `push(array, value)`
  Returns a new array with `value` pushed onto the end of `array`.
- `pop(array)`
  Returns the last value of the `array` or `null` if empty.
- `exit([status])`
  Exits the program immediately with the optional `status` or `0`.
- `assert(expr, [msg])`
  Exits the program immediately with a non-zero status if `expr` is `false`
  optionally displaying `msg` to standard error.
- `bool(value)`
  Converts `value` to a `bool`. If `value` is `bool` returns the value directly.
  Returns `true` for non-zero `int`(s), `false` otherwise. Returns `true` for
  non-empty `str`, `array` and `hash` values. Returns `true` for all other
  values except `null` which always returns `false`.
- `int(value)`
  Converts decimal `value` `str` to `int`. If `value` is invalid returns `null.
  If `value` is an `int` returns its value directly.
- `str(value)`
  Returns the string representation of `value`: `null` for null,
  `true` or `false` for `bool`, decimal for `int` (eg: `1234`),
  the string itself for `str` (not quoted),
  the Monkey representation for array and hash (eg: `[1, 2]` and `{"a": 1}`
  with keys sorted), and something like `<fn name(...) at 0x...>` for functions..
- `type(value)`
  Returns a `str` denoting the type of value: `nil`, `bool`, `int`, `str`, `array`, `hash`, or `fn`.
- `args()`
  Returns an array of command-line options passed to the program.
- `lower(str)`
  Returns a lowercased version of `str`.
- `upper(str)`
  Returns an uppercased version of `str`.
- `join(array, sep)`
  Concatenates `str`s in `array` to form a single `str`, with the separator `str` between each element.
- `split(str[, sep])`
  Splits the `str` using given separator `sep`, and returns the parts (excluding the separator) as an `array`. If `sep` is not given or `null`, it splits on whitespace.
- `find(haystack, needle)
  Returns the index of `needle` `str` in `haystack` `str`,
  or the index of `needle` element in `haystack` array.
  Returns -1 if not found.
- `readfile(filename)`
  Reads the contents of the file `filename` and returns it as a `str`.
- `writefile(filename, data)`
  Writes `data` to a file `filename`.
- `abs(n)`
  Returns the absolute value of the `n`.
- pow(x, y)`
  Returns `x` to the power of `y` as `int`(s).
- `divmod(a, b)`
  Returns an array containing the  quotient and remainder of `a` and `b` as `int`(s).
  Equivilent to `[a / b, b % b]`.
- `bin(n`)
  Returns the binary representation of `n` as a `str`.
- `hex(n)`
  Returns the hexidecimal representation of `n` as a `str`.
- `oct(n)`
  Returns the octal representation of `n` as a `str`.
- `ord(c)`
  Returns the ordincal value of the character `c` as an `int`.
- `chr(n)`
  Returns the character value of `n` as a `str`.
- `hash(any)`
  Returns the hash value of `any` as an `int`.
- `id(any)`
  Returns the identity of `any` as an `int`.
- `min(array)`
  Returns the minimum value of elements in `array`.
- `max(array)`
  Returns the maximum value of elements in `array`.
- `sorted(array)`
  Sorts the `array` using a stable sort, and returns  a new `array`..
  Elements in the `array` must be orderable with `<` (`int`, `str`, or `array` of those).
- `reversed(array)`
  Reverses the array `array` and returns a new `array`.
- `open(filename[, mode])`
- `write(fd, data)`
  Writes `str` `data` to the open file descriptor given by `int` `fd`.
- `read(fd, [n])`
  Reads from the file descriptor `fd` (`int`) optinoally up to `n` (`int`)
  bytes and returns the read data as a `str`.
- `close(fd)`
  Closes the open file descriptor given by `fd` (`int`).
- `seek(fd, offset[, whence])`
  Seeks the file descriptor `fd` (`int`) to the `offset` (`int`). The optional
  `whence` (`int`) determins whether to seek from the beginning of the file (`0`),
  relativie to the current offset (`1`) or the end of the file (`2`).
- `socket(type)`
- `bind(fd, address)`
- `listen(fd, backlog)`
- `accept(fd)`
- `connect(fd, address)`

### Objects

```#!sh
>> Person := fn(name, age) { self := {} self.name = name self.age = age self.str = fn() { return self.name + ", aged " + str(self.age) } return self }
>> p := Person("John", 35)
>> p.str()
"John, aged 35"
```

### Modules

Monkey supports modules. Modules are just like other Monkey source files
with the extension `.monkey`.  Modules are searched for by `SearchPaths`
which can be controlled by the environment `MONKEYPATH`. By default this is
always the current directory.

To import a module:

```
>> foo := import("foo")
>> foo.A
5
>> foo.Sum(2, 3)
5
```

## License

This work is licensed under the terms of the MIT License.
