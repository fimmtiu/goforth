# goforth

A tiny Forth compiler & virtual machine written in about a day.

## Usage

```
$ go build
$ echo ': foo 1 2 + ; foo .' | ./goforth
3
```

It'll read Forth code from standard input. It's about as minimal a feature set as you can get: it can do `if else then`, `+`, `.`, user-defined words, and not much else. My goal was to get it to a point where it could run FizzBuzz.

## Notes

You can make much smaller and more elegant Forth interpreters, but the goal of this project was to muck about with compilers and virtual machines. It's a little stack-based virtual machine with a 32-bit instruction set. Go is not a great implementation language for this sort of thing, and I frequently found myself wishing I'd done this in C instead, but that's what I get for wanting to practice Go.

It's a bit of a dog's breakfast, code-wise, given the haste in which it was written. Please be kind.
