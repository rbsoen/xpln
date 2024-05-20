# xpln: Explain your code

`xpln` is a literate programming tool written in Go. It can be used for introducing concepts, creating programming tutorials, or (in general) to structure a piece of software as if it were code.

## Instructions

`xpln` works with .md files with the following format:

1. A named code block is started with a single line of the following format: 3 backticks; the language name; the code block's name. To end it, write a line consisting of only backticks (min. 3).
2. Code blocks can reference another code block by surrounding its name with `@{` and `}`. The reference can be indented, in which case the referenced block will have the same indent on every line. The line in which this reference is contained must only have said reference in it.

`xpln` can generate usable source code and web pages.

For single-HTML web pages:
```
xpln weave source1.md source2.md > index.html
```

For source code (to `sourcedir`):
```
xpln tangle sourcedir source1.md source2.mid
```

Note: The files output will only be code blocks that start with a /. Its file name is everything after the /, and subdirectories will also be creeated with it.
