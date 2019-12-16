// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package compile compiles Suneido source code
to byte code to be interpreted by the runtime package.

It uses a recursive descent parser that produces an AST
that is then converted to byte code by codegen

Compiling constants (other than functions)
produces values directly without going through an AST.

Expression parsing is shared with database queries.
*/
package compile
