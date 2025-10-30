### Generator

Abstract base class for generator formats. Defines `Generator?()` to return true.

Derived classes define:
`Header( ) => format or false`
: Returns a format item that will be printed at the beginning of the generated output, and also at the top of each additional page of generated output. Can return false if no header is desired. The default definition returns false.

`Next( ) => format or false`
: Returns the next format item or false if all output has been generated.

`More()`
: Instead of supplying a Next function, generators can supply a More function that calls .Output (Generator.Output) to place items in the output queue. Note: If a call to More does not produce any output, this is taken as the end of generation.

`Output(format)`
: Used by More() to place format items in the output queue. These items will be produced by the generator in the order they are output. (i.e. first in first out).