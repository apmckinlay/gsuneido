// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// TAGS: gsuneido
/*
standardizes:
	:name
	commas between arguments and members
	quotes and escaping as per Display
	.mem instead of this.mem
	x.y instead of x["y"]
	#(foo) instead of #("foo")
	forever instead of for(;;) or while(true)
	++x and --x when value not used, instead of x++ or x--
	[...] instead of Object(...) and Record(...), except for Object with only named
*/
//TODO don't add comma in Assert(x, is: 1)
//MAYBE keep zero indent (debugging)
//TODO multi-line strings
class
	{
	prevPos: 0
	output: ""
	indent: 0

	CallClass(src)
		{
		return (new this).Process(src)
		}
	Process(.src)
		{
		ob = Suneido.Parse(src, true)
		ast = ob[0]
		.comments = ob[1]
		_constant = true
		.fmt(ast)
		if .output[-1] isnt '\n'
			.nl()
		return .output
		}
	fmt(node)
		{
		if node is false
			return
		_node = node
		c = .comments[node]
		if c isnt false and c[0] isnt #()
			{
.trace(node.type.Upper(), before: Display(c[0]))
			.printExtras(c[0])
			}
.trace(node.type.Upper(), node)
		this[node.type](node)
		if c isnt false and c[1] isnt #()
			{
.trace(node.type.Upper(), after: Display(c[1]))
			.printExtras(c[1])
			}
		}
	printExtras(extra)
		{
		for s in extra
			{
			if s is '\n'
				.extra_nl()
			else if s[0] is '/'
				.printExtra(s)
			}
		}

	Object(node) // constant
		{
		.members(node, "()")
		}
	Record(node) // constant
		{
		.members(node, "{}")
		}
	members(node, delims)
		{
		if not _constant or .output is ""
			.print('#', nosuf:)
		.print(delims[0], nosuf:)
		_constant = true
		.sep = ""
		for (i = 0; false isnt kv = node[i]; ++i)
			{
			// kv is either [k] for unnamed, or [k,v] for named
			if kv.Size() is 1
				.constant(kv[0])
			else
				{
				.constant(kv[1])
				.print(':', nopre:)
				.memval(kv[0])
				}
			.sep = ', '
			}
		.print(delims[1], nopre:)
		}
	constant(x)
		{
		if Type(x) is #AstNode
			.fmt(x)
		else if _constant and String?(x) and x.Identifier?() and
			x not in (#true, #false, #function, #class, #dll, #struct, #callback)
			.print(x)
		else
			.print(Display(x))
		}
	memval(x)
		{
		if Type(x) is #AstNode and x.type is #Constant and x.value is true
			{
			// in k: true suppress the true, but print its comments
			c = .comments[x]
			if c isnt false and c[0] isnt #()
				.printExtras(c[0])
			if c isnt false and c[1] isnt #()
				.printExtras(c[1])
			}
		else if x isnt true
			.constant(x)
		}

	Class(node)
		{
		.print(node.base $ ' {')
		.funcExtras(node)
		.nl()
		++.indent
		_constant = true
		for (i = 0; false isnt kv = node[i]; ++i)
			{
			if Type(kv[0]) is #AstNode
				{
				if kv[0].type is #Function
					{
					.print(kv[1], nosuf:)
					.func(kv[0])
					}
				else
					{
					.constant(kv[1])
					.print(':', nopre:)
					.fmt(kv[0])
					}
				.nl()
				}
			else
				{
				.constant(kv[1])
				.println(': ' $ Display(kv[0]), nopre:)
				}
			}
		--.indent
		.print('}')
		}
//	mem(m)
//		{
//		if Type(m) is 'AstNode'
//			{
//			_constant = true
//			.fmt(m)
//			.print(':', nopre:)
//			}
//		else
//			.print((m.Identifier?() ? m : Display(m)) $ ':')
//		}

	Function(node)
		{
		.print("function", nosuf:)
		.func(node)
		}
	func(node)
		{
		.params(node.params, "()")
		.print(" {", nopre:)
		.funcExtras(node)
		.funcBody(node)
		}
	funcExtras(node)
		{
		if false isnt c = .comments[node]
			{
.trace(extras: c[1])
			prev = false
			for s in c[1]
				{
				if s[0] isnt '{'
					break
				if s isnt "{\n"
					{
					// suppress extra newlines unless followed by a comment
					if prev is "{\n"
						.extra_nl()
					.printExtra(s[1..])
					}
				prev = s
				}
			}
		}
	Block(node)
		{
		.print("{", nosuf:)
		params = node.params
		if params.size > 1 or (params.size is 1 and params[0].name isnt "it")
			.params(params, "||")
		else
			.print("")
		.funcBody(node)
		}
	params(params, delims)
		{
		_constant = false // want quotes and #
		_singleLine = true
		.print(delims[0], nosuf:)
		for (i = 0; false isnt p = params[i]; ++i)
			{
			.fmt(p)
			.sep = ", "
			}
		.print(delims[1], nopre:)
		}
	Param(node)
		{
		.print(node.name)
		if false isnt defval = node.defval
			{
			.print('=')
			.constant(defval)
			}
		}
	funcBody(node, _singleLine = false)
		{
		_constant = false
		_singleLine = true
		if singleLine and .okToSingleLine(node)
			{
			for (i = 0; false isnt stmt = node[i]; ++i)
				{
				if stmt.type is #Compound and stmt.size is 0
					{
					if i is node.size - 1
						.print(';;', nopre:)
					break
					}
				.fmt(stmt)
				.sep = "; "
				}
			.print(" }", nopre:)
			}
		else
			{
			.nl()
			.statements(node)
			.print("}", nopre:)
			}
		}
	okToSingleLine(node)
		{
		if false isnt c = .comments[node]
			if c[1].Size() > 0 and c[1][0][0] is '{'
				return false
		n = node.size
		if n is 0
			return true
		if node[n-1].type is #Compound and node[n-1].size is 0
			--n
		for (i = 0; i < n; ++i)
			if node[i].type not in (#Expr ,#Return, #Throw)
				return false
		return true
		}

	Expr(node) // statement
		{
		.expr(node.expr)
		}
	expr(node)
		{
		// convert PostInc/Dec to pre Inc/Dec when result not used
		if node.type is #Unary and node.op in (#PostInc, #PostDec)
			.fmt(Object(type: #Unary, expr: node.expr,
				op: node.op.RemovePrefix("Post")))
		else
			.fmt(node)
		}

	Constant(node, _constant = false)
		{
		x = node.value
		if Type(x) is #AstNode
			.fmt(x)
		else if node.symbol
			.print('#' $ x)
		else if x is 1/0
			.print("1/0")
		else if x is -1/0
			.print("-1/0")
		else if constant and String?(x) and x.Identifier?() and
			x not in (#true, #false, #function, #class, #dll, #struct, #callback)
			.print(x)
		else
			.print(Display(x))
		}

	Ident(node)
		{
		.print(node.name)
		}

	Unary(node)
		{
		if node.op is #Not and node.expr.type is #In
			return .In(node.expr, "not in")
		op = #(Add: '+', Sub: '-', Not: 'not', BitNot: '~', Inc: '++', PostInc: '++',
			Dec: '--', PostDec: '--', Div: '1/')
		if node.op is #LParen
			.print("(", nosuf:)
		else if node.op not in ("PostInc", "PostDec")
			.print(op[node.op], nosuf: node.op isnt #Not)
		.fmt(node.expr)
		if node.op is #LParen
			.print(")", nopre:)
		else if node.op in ("PostInc", "PostDec")
			.print(op[node.op], nopre:)
		}

	Binary(node)
		{
		op = #(Eq: '=', AddEq: '+=', SubEq: '-=', CatEq: '$=', MulEq: '*=', DivEq: '/=',
			ModEq: '%=', LShiftEq: '<<=', RShiftEq: '>>=', BitOrEq: '|=', BitAndEq: '&=',
			BitXorEq: '^=', Is: 'is', Isnt: 'isnt', Match: '=~', MatchNot: '!~', Mod: '%',
			LShift: '<<', RShift: '>>', Lt: '<', Lte: '<=', Gt: '>', Gte: '>=')
		.fmt(node.lhs)
		.print(op[node.op])
		.fmt(node.rhs)
		}

	Nary(node)
		{
		ops = #(And: 'and', Or: 'or', Add: '+', Cat: '$', Mul: '*',
			BitOr: '|', BitAnd: '&', BitXor: '^')
		op = ops[node.op]
		.fmt(node[0])
		for (i = 1; i < node.size; ++i)
			{
			if op is '+' and node[i].type is #Unary and node[i].op is #Sub
				{
				.print('-')
				.fmt(node[i].expr)
				}
			else if op is '+' and node[i].type is #Constant and node[i].value < 0
				{
				.print('-')
				.fmt([type: #Constant, value: -node[i].value, symbol: false])
				}
			else if op is '*' and node[i].type is #Unary and node[i].op is #Div
				{
				.print('/')
				.fmt(node[i].expr)
				}
			else
				{
				.print(op)
				.fmt(node[i])
				}
			}
		}

	Trinary(node)
		{
		.fmt(node.cond)
		.print("?")
		.fmt(node.t)
		.print(":")
		.fmt(node.f)
		}

	Mem(node)
		{
		if node.mem.type is #Constant and
			String?(node.mem.value) and node.mem.value.Identifier?()
			{
			if node.expr.type isnt #Ident or node.expr.name isnt "this"
				{
				.fmt(node.expr)
				.sep = ''
				}
			.print("." $ node.mem.value)
			}
		else
			{
			.fmt(node.expr)
			.print("[", nopre:, nosuf:)
			.fmt(node.mem)
			.print("]", nopre:)
			}
		}

	RangeTo(node)
		{
		.fmt(node.expr)
		.print('[', nopre:, nosuf:)
		.fmt(node.from)
		simple = .simple(node.from) and .simple(node.to)
		.print("..", nopre: simple, nosuf: simple)
		.fmt(node.to)
		.print("]", nopre:)
		}
	RangeLen(node)
		{
		.fmt(node.expr)
		.print('[', nopre:, nosuf:)
		.fmt(node.from)
		simple = .simple(node.from) and .simple(node.len)
		.print("::", nopre: simple, nosuf: simple)
		.fmt(node.len)
		.print("]", nopre:)
		}
	simple(expr)
		{
		return expr is false or expr.type in (#Constant, #Ident)
		}

	In(node, op = "in")
		{
		.fmt(node.expr)
		.print(op $ " (", nosuf:)
		for (i = 0; false isnt arg = node[i]; ++i)
			{
			.fmt(arg)
			.sep = ", "
			}
		.print(')', nopre:)
		}

	Call(node, _cond = false)
		{
		if .useBrackets(node)
			{
			.print('[', nosuf:)
			delim = ']'
			}
		else
			{
			.fmt(node.func)
			.print("(", nopre:, nosuf:)
			delim = ')'
			}
		for (i = 0; false isnt arg = node[i]; ++i)
			{
			if not cond and i is node.size - 1 and
				arg.name is "block" and arg.expr.type is "Block"
				{
				//TODO comments
				.print(delim, nopre:)
				_singleLine = false // block outside parens is multi-
				.Block(arg.expr)
				return
				}
			if i > 0
				.print(", ", nopre:)
			.fmt(arg)
			}
		.print(delim, nopre:)
		}
	useBrackets(node)
		{
		if node.func.type isnt #Ident
			return false
		if node.func.name is #Record
			return true
		if node.func.name isnt #Object
			return false
		arg = node[0]
		return arg isnt false and arg.name is false
		}
	Argument(arg, _cond = false)
		{
		if arg.name isnt false and arg.expr.type is #Ident and
			arg.expr.name.LocalName?() and arg.expr.name is arg.name
			.print(':' $ arg.name)
		else
			{
			if arg.name isnt false
				if arg.name in ('@', "@+1")
					.print(arg.name, nosuf:)
				else if String?(arg.name) and arg.name.Identifier?()
					.print(arg.name $ ':')
				else
					.print(Display(arg.name) $ ':')
			if arg.name is false or
				arg.expr.type isnt #Constant or arg.expr.value isnt true
				.fmt(arg.expr)
			}
	}

	Compound(node)
		{
		.println('{')
		for (i = 0; false isnt stmt = node[i]; ++i)
			{
			.fmt(stmt)
			.nl()
			}
		.println('}')
		}

	Return(node)
		{
		.print("return")
		.fmt(node.expr)
		}

	Throw(node)
		{
		.print("throw")
		.fmt(node.expr)
		}

	TryCatch(node)
		{
		.print("try")
		.body(node.try)
		if false isnt c = node.catch
			{
			if .output[-1] isnt '}'
				.nl()
			.print("catch")
			if c.Size() > 1
				{
				.print('(' $ c[1], nosuf:)
				if c.Size() > 2
					.print(", " $ Display(c[2]))
				.print(')', nopre:)
				}
			.body(c[0])
			}
		}

	If(node)
		{
		forever // loop over else-if
			{
			.print("if")
			.cond(node.cond)
			.body(node.t, needBraces: .needBraces(node))
			if node.f isnt false
				{ // else
				if .output[-1] isnt '}'
					.nl()
				.print("else")
				if node.f.type is #If
					{
					node = node.f
					continue
					}
				.body(node.f)
				}
			break
			}
		}
	cond(expr)
		{
		_cond = true
		.fmt(expr)
		}
	needBraces(node)
		{
		if node.f is false
			return false
		t = .unwrap(node.t)
		return t.type is #If and t.f is false
		}
	unwrap(node)
		{
		while node.type is #Compound and node.size is 1
			node = node[0]
		return node
		}

	Switch(node)
		{
		.print("switch")
		if node.expr isnt false and
			(node.expr.type isnt #Constant or node.expr.value isnt true)
			.cond(node.expr)
		.print('{')
		.funcExtras(node)
		.nl()
		for (i = 0; false isnt c = node[i]; ++i)
			{
			.print("case")
			for (j = 0; j < c.Size()-1; ++j, .sep = ", ")
				.fmt(c[j])
			.print(":", nopre:)
			.nl()
			.statements(c.Last())
			}
		if node.def isnt false
			{
			.print("default:")
			.nl()
			.statements(node.def)
			}
		.print("}")
		}

	Forever(node)
		{
		.print("forever")
		.body(node.body)
		}
	While(node)
		{
		if node.cond.type is #Constant and node.cond.value is true
			return .Forever(node)
		.print("while")
		.cond(node.cond)
		.body(node.body)
		}
	DoWhile(node)
		{
		.print("do")
		body = node.body
		if body.type isnt #Compound
			body = [body, false, type: #Compound, size: 1]
		.body(body, needBraces:)
		++.indent
		.print("while")
		.cond(node.cond)
		--.indent
		}
	ForIn(node)
		{
		.print("for", node.var, "in")
		.cond(node.expr)
		.body(node.body)
		}
	For(node)
		{
		if node.init.Empty?() and node.cond is false and node.inc.Empty?()
			return .Forever(node)
		.print("for (", nosuf:)
		for e in node.init
			{
			.expr(e)
			.sep = ", "
			}
		.print(';', nopre:)
		.cond(node.cond)
		.print(';', nopre:)
		for e in node.inc
			{
			.expr(e)
			.sep = ", "
			}
		.print(')', nopre:)
		.body(node.body)
		}
	body(node, needBraces = false)
		{
		while not needBraces and node.type is #Compound and node.size is 1
			node = node[0]
		if node.type is #Compound
			{
			.println("{")
			.statements(node)
			.print("}")
			}
		else
			{
			.nl()
			++.indent
			.fmt(node)
			--.indent
			}
		}
	statements(node)
		{
		++.indent
		for (i = 0; false isnt stmt = node[i]; ++i)
			{
			.fmt(stmt)
			if i is node.size - 2 and node[i+1].type is #Compound and node[i+1].size is 0
				{
				.print(';;', nopre:)
				++i
				}
			.nl()
			}
		--.indent
		}
	Break(unused)
		{
		.print("break")
		}
	Continue(unused)
		{
		.print("continue")
		}

	sep: ""
	extraLine: true // whether the current line is all comments or blank
	println(@args)
		{
		.print(@args)
		.nl()
		}
	print(@args)
		{
		.print_(args)
		.extraLine = false
		}
	printExtra(@args)
		{
		if args[0] is '/*unused*/' and _node.type is #Param and _node.name is "unused"
			return // suppress /*unused*/ if parameter name is "unused"
		.print_(args)
		}
	print_(args)
		{
.trace(print: args)
		if .output.Size() is .last_extra_nl
			{
			.output $= '\t'.Repeat(.indent + (.extraLine ? 0 : 1))
			.sep = ""
			.extraLine = true
			}
		else if .output[-1] is '\n'
			{
			.output $= '\t'.Repeat(.indent)
			.sep = ""
			.extraLine = true
			}
		if args.Extract(#nopre, false) is true
			.sep = ""
		nosuf = args.Extract(#nosuf, false)
		.prev = .output.Size()
		for arg in args
			{
			if .output[-1] isnt .sep
				.output $= .sep
			.output $= arg
			.sep = " "
			}
		if nosuf
			.sep = ""
		}
	nl()
		{
.trace('auto nl')
		if .output.Size() is .last_extra_nl
			.last_extra_nl = false
		else
			.output $= '\n'
		}
	last_extra_nl: false
	extra_nl()
		{
.trace('extra nl')
		if .output.Suffix?("\n\n") // don't allow more than one blank line
			return
		// merge with automatic newline
		if .output[-1] isnt '\n' or .output.Size() is .last_extra_nl
			{
			.output $= '\n'
			.last_extra_nl = .output.Size()
			}
		}
trace(@args)
{
//Print(@args)
}
	}