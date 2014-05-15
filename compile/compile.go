/*
Package compile implements compiling Suneido source code
to byte code to be interpreted by the interp package.

It uses a recursive descent parser that produces an AST
that is then converted to byte code by codegen

Compiling constants (other than functions)
produces values directly without going through an AST.

Expression parsing is shared with database queries.
*/
package compile
