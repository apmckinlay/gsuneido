<div style="float:right"><span class="builtin">Builtin</span></div>

#### Suneido.Parse

``` suneido
(string) => ast
```

Suneido.Parse parses a string containing a function and returns an abstract syntax tree (AST). It uses the parser from Suneido's compiler. The returned AST is a "view" of the internal AST. To parse a query, see Query.Parse

An ast will display as a string in an internal format.

``` suneido
Suneido.Parse("function () { f(); g() }")
=> "Function(
        Call(f)
        Call(g))"
```

Nodes have "properties". Each node has a **type** (a string). For example:

``` suneido
ast = Suneido.Parse("function(){}")
ast.type => "Function"
```

Nodes with a variable number of child nodes (Function, Nary, In, Call, Compound, Switch) have a **size** property and can be indexed. For example:

``` suneido
ast = Suneido.Parse("function(){ f(); g() }")
ast.size    => 2
ast[1]      => "Call(g)"
```

All nodes have **pos** and **end** properties. If available, pos...end will be the position of the node in the source code, otherwise they will be false.

Nodes also have a **children** property that can be used to generically traverse the AST and visit every node.

``` suneido
AstTraverse
function (node)
	{
	if Type(node) is 'AstNode'
		for (i = 0; false isnt c = node.children[i]; ++i)
			AstTraverse(c)
	}
```

#### Constant Node Types and Properties

#### Object, Record, Class

base => string

size => number

[i] => Member

#### Member

value => value

named => true or false

key => value (only if named is true)

Note: Class members are always named (have keys).

#### Function

params => Params

size => number

[i] => Statement

#### Params
size => number
[i] => Param
#### Param

name => string

hasdef => true or false

defval => value

unused => true or false

#### Expression Node Types and Properties

#### Constant

value

symbol => true/false

Note: symbol distinguishes "a string" from #symbol

#### Ident

name => string

#### Unary

op => string

Add, Sub, Not, BitNot, Inc, PostInc, Dec, PostDec, Div, LParen

Note: unary divide is reciprocal i.e. 1/n

Note: unary LParen is e.g (.fn)()

expr

#### Binary

lhs => expression

op => string

Eq, AddEq, SubEq, CatEq, MulEq, DivEq, ModEq,
LShiftEq, RShiftEq, BitOrEq, BitAndEq, BitXorEq,
Is, Isnt, Match, MatchNot, Mod,
LShift, RShift, Lt, Lte, Gt, Gte

rhs => expression

#### Nary

op => string

And, Or, Add, Cat, Mul, BitOr, BitAnd, BitXor

Note: subtract is converted to add of unary Sub
and divide is converted to multiply by unary Div (reciprocal)

size => number

[i] => expression

#### Trinary ?:

cond => expression

t => expression

f => expression

#### Mem

expr

mem => expression

#### RangeTo

expr

from => number or false

to => number or false

#### RangeLen

expr

from => number or false

len => number or false

#### InRange

InRange is internally produced by combining e.g. x > 0 and x <= 10 => InRange(0, 10)

expr

from => constant

to => constant

#### In

expr => expression

size => number

[i] => expression

#### Call

func => expression

size => number

[i] => Argument

#### Argument

name => string or false

expr => expression

#### Statement Nodes and Properties

#### ExprStmt

expr

#### Compound

size => number

[i] => statement

#### If

cond => expression

t => statement

f => statement

#### Switch

expr

size => number

[i] => Case

def => compound

#### Case

[i] => expression

body => Compound

#### Return

expr => expression or false

throw => true if return throw, otherwise false

#### Throw

expr

#### TryCatch

try => statement

catch => statement

catchvar => Ident or false

catchpat => string or false

#### Forever

body => statement

#### ForIn

var => string, will be empty with `for ..n`

expr

expr2 - used by `for i in j..k` or `for ..n`

body => statement

#### For

init => #(expressions)

cond => expression

inc => #(expressions)

body => statement

#### While

cond => expression

body => statement

#### DoWhile

cond => expression

body => statement

#### Break

#### Continue