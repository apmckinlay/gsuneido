// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
/*
standardizes:
	:name
	commas between arguments and members
	.mem instead of this.mem
	x.y instead of x["y"]
	#(foo) instead of #("foo")
	forever instead of for(;;) or while(true)
	++x and --x when value not used, instead of x++ or x--
	[...] instead of Object(...) and Record(...), except for Object with only named
*/
//MAYBE normalize comments // at end of line /**/ otherwise
//MAYBE build current line separately so we can examine/modify, when final add to output
class
	{
	output: ""
	indent: 0

	CallClass(src)
		{
		return (new this).Process(src)
		}
	Process(.src)
		{
		maxpos = 9999999
		ast = Suneido.Parse(src)
		.scan = Scanner(src)
		_constant = false
		_parentEnd = maxpos
		.softLineBreak = .extraLineBreak = false
		.fmt(ast)
		if .output[-1] isnt '\n' //FIXME will probably flatten
			.nl()
		.trailing(maxpos)
		return .output
		}
	fmt(node, _parentEnd, allowBlankLineAfter = false, noBreakBefore = false)
		{
		if node is false
			return
		_node = node // used by printExtra for /*unused*/
.trace("(((", node.type.Upper(), node)
.trace(:allowBlankLineAfter, :noBreakBefore)
		if false isnt pos = node.pos
			{ // leading
end = node.end
.trace(pos, end, Display(.src[pos..end]), parentEnd)
			if node.end isnt 0
				_parentEnd = node.end // constrain child node extra processing
			.leading(node, pos, :noBreakBefore)
			}
		this[node.type](node)
		if false isnt pos and node.end not in (0, false)
			.trailing(node.end, parentEnd, :allowBlankLineAfter)
.trace(")))", node.type.Upper(), node)
		}

	// leading prints comments up to node.pos
	leading(node, pos, noBreakBefore = false)
		{
		while .scan.Position() < pos and .scan isnt .scan.Next2()
			{
.trace(.scan.Position(), leading_scan: Display(.scan.Text()))
			switch .scan.Type()
				{
			case #COMMENT:
				.printExtra(.scan.Text())
			case #NEWLINE:
				if not noBreakBefore
					{
					if node.type in (#Ident, #Constant, #Unary, #Mem)
						.softLineBreak()
					.extraLineBreak()
					}
			default:
				//ignore
				}
			}
		}

	// trailing prints comments before pos
	// and after pos up to a token that's not a comment or whitespace
	trailing(pos, parentEnd = 9999999, allowBlankLineAfter = false)
		{
.trace("trailing", pos, parentEnd, Display(.src[.scan.Position()..Min(pos, parentEnd)]))
		while .scan.Position() < parentEnd and .scan isnt .scan.Next2()
			{
.trace(.scan.Position(), trailing_scan: Display(.scan.Text()))
			switch .scan.Type()
				{
			case #COMMENT:
				.printExtra(.scan.Text())
			case #WHITESPACE:
				// ignore
			case #NEWLINE:
				.extraLineBreak()
				if allowBlankLineAfter and .blankLine?(.scan.Text())
					.blankLine()
				if .scan.Position() >= pos
					return
			default:
				if .scan.Position() > pos
					return
				}
			}
		}
	blankLine?(s)
		{
		return s.Find('\n', s.Find('\n')+1) < s.Size()
		}

	Object(node) // constant
		{
		.members(node, "()")
		}
	Record(node) // constant
		{
		.members(node, "{}")
		}
	nest: 0 // used by members and Call
	members(node, delims)
		{
		if not _constant or .output is ""
			.print('#', nosuf:)
		.print(delims[0], nosuf:)
		_inclass = false
		_constant = true // if non-constant it would be arguments
		.sep = ""
		oldnest = .nest
		for (i = 0; false isnt m = node[i]; ++i)
			{
			if i is 1
				++.nest
			_last = false is node[i+1] // used by objectMember to know if to add a comma
			.softLineBreak()
			.fmt(m, allowBlankLineAfter: false isnt node[i+1]) // => Member
			}
		.nest = oldnest
		.softLineBreak()
		.print(delims[1], nopre:)
		}
	Member(m, _inclass = false) // container member
		{
		if inclass
			.classMember(m)
		else
			.objectMember(m)
		}
	objectMember(m, _last = false)
		{
		if not m.named
			.constant(m.value)
		else
			{
			.constant(m.key)
			.print(':', nopre:)
			if m.value isnt true
				.constant(m.value)
			}
		if not last
			.print(',', nopre:)
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

	Class(node)
		{
		.print(node.base)
		.trailing(node.pos1)
		.nl()
		++.indent
		.print('{')
		.trailing(node.pos2)
		.nl()
		_inclass = true
		_constant = true
		for (i = 0; false isnt m = node[i]; ++i)
			{
			.fmt(m, allowBlankLineAfter: false isnt node[i+1])
			.nl()
			}
		.print('}')
		--.indent
		}
	classMember(m)
		{
		if Type(m.value) is #AstNode
			{
			if m.value.type is #Function
				{
				.print(m.key, nosuf:)
				.func(m.value)
				}
			else
				{
				.constant(m.key)
				.print(':', nopre:)
				.fmt(m.value)
				}
			}
		else
			{
			.constant(m.key)
			.print(':', nopre:)
			.print(Display(m.value))
			}
		}

	Function(node)
		{
		.print("function", nosuf:)
		.func(node)
		}
	func(node) // used by Function and classMember
		{
		.params(node.params, "()")
		.trailing(node.pos1)
		.funcBody(node)
			{ .trailing(node.pos2) }
		}
	Block(node)
		{
		.funcBody(node)
			{
			params = node.params
			if params.size > 1 or (params.size is 1 and params[0].name isnt "it")
				.params(params, "||")
			}
		}
	params(params, delims)
		{
		_constant = false // want quotes and #
		_singleLine = true
		.print(delims[0], nopre:, nosuf:)
		for (i = 0; false isnt p = params[i]; ++i)
			{
			if i > 0
				{
				.print(',', nopre:)
				.sep = ' '
				.softLineBreak()
				}
			.fmt(p)
			}
		.print(delims[1], nopre:)
		}
	Param(node)
		{
		.print(node.name)
		if node.hasdef
			{
			.trailing(node.pos + node.name.Size())
			.print('=')
			.constant(node.defval)
			}
		}
	funcBody(node, _singleLine = false, block = function(){})
		{
		_constant = false
		_singleLine = true
		if singleLine and .okToSingleLine(node)
			{
			.print("{")
			block()
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
			++.indent
			.print("{")
			block()
			.nl()
			.statements2(node)
			.leading(node, node.end)
			.print("}", nopre:)
			--.indent
			}
		}
	okToSingleLine(node)
		{
		if node.pos isnt false and .src[node.pos..node.end].Has?('\n')
			return false
		return .okToSingleLine2(node)
		}
	okToSingleLine2(node)
		{
		if Type(node) isnt 'AstNode'
			return true
		if node.type in
			(#If, #Switch, #TryCatch, #For, #ForIn, #Forever, #While, #DoWhile)
			return false
		children = node.children
		for (i = 0; false isnt c = children[i]; ++i)
			if not .okToSingleLine2(c)
				return false
		return true
		}

	ExprStmt(node)
		{
		.expr(node.expr)
		}
	expr(node)
		{
		// convert PostInc/Dec to pre Inc/Dec when result not used
		if node.type is #Unary and node.op in (#PostInc, #PostDec)
			.fmt(Object(type: #Unary, expr: node.expr,
				op: node.op.RemovePrefix("Post"), pos: node.pos, end: node.end))
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
		else if constant and String?(x) and x.Identifier?() and
			x not in (#true, #false, #function, #class, #dll, #struct, #callback)
			.print(x)
		else if node.pos not in (0, false) and (String?(x) or Number?(x))
			.print(.src[node.pos..node.end])
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
			else if op is '+' and node[i].type is #Constant and
				Number?(node[i].value) and node[i].value < 0
				{
				.print('-')
				.fmt(Object(type: #Constant, value: -node[i].value, symbol: false,
					pos: node[i].pos, end: node[i].end))
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
		.softLineBreak()
		.print("?")
		.fmt(node.t, noBreakBefore:)
		.softLineBreak()
		.print(":")
		.fmt(node.f, noBreakBefore:)
		}

	Mem(node) // a.b or .b or a[b]
		{
		if node.mem.type is #Constant and
			String?(node.mem.value) and node.mem.value.Identifier?()
			{
			if node.expr.type isnt #Ident or node.expr.name isnt "this"
				{
				.fmt(node.expr)
				.sep = ''
				}
			else
				.leading(node, node.dotpos)
			.print(".", nosuf:)
			.softLineBreak()
			.trailing(node.dotpos)
			.print(node.mem.value)
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
		if .new_this(node)
			{
			.print("new this")
			if node.size is 0
				delim = ""
			else
				{
				.print('(', nopre:, nosuf:)
				delim = ')'
				}
			}
		else if .useBrackets(node)
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
		oldnest = .nest
		last = node.size - 1
		if last >= 0 and not cond
			{
			arg = node[last]
			if arg.name is "block" and arg.expr.type is "Block"
				--last
			}
		for (i = 0; false isnt arg = node[i]; ++i)
			{
			if i is 1
				++.nest
			if i is last + 1
				{
				.print(delim, nopre:)
				.nest = oldnest
				_singleLine = false // block outside parens is multi-line
				.Block(arg.expr)
				return
				}
			if i > 0
				.sep = ' '
			.softLineBreak()
			_argsep = i is last ? "" : ','
			.fmt(arg) // => Argument
			}
		.softLineBreak()
		.print(delim, nopre:)
		.nest = oldnest
		.trailing(node.end, _parentEnd, allowBlankLineAfter:)
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
	new_this(node)
		{
		if #Mem is (fn = node.func).type and
			#Ident is (expr = fn.expr).type and expr.name is 'this' and
			#Constant is (meth = fn.mem).type and meth.value is '*new*'
			{
			return true
			}
		return false
		}
	Argument(arg, _argsep = ',')
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
		.output $= argsep
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
		.trailing(node.pos + 3) /*= "try".Size() */
		.body(node.try)
		if false isnt c = node.catch
			{
			.nl()
			.print("catch")
			if false isnt var = node.catchvar
				{
				.print('(', nopre:, nosuf:)
				.fmt(var)
				if false isnt pat = node.catchpat
					.print(", " $ Display(pat), nopre:, nosuf:)
				.print(')', nopre:)
				}
			.trailing(node.catchend)
			.body(c)
			}
		}

	If(node)
		{
		forever // loop over else-if
			{
			.print("if")
			.cond(node.cond)
			.body(node.t)
			if node.f isnt false
				{ // else
				.nl()
				.print("else")
				.trailing(node.elseend)
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
		_cond = true // used in Call to avoid block after parens
//		if expr.type is #Unary and expr.op is #LParen
//			expr = expr.expr
		.fmt(expr)
		}

	Switch(node)
		{
		.print("switch")
		if node.expr isnt false and
			(node.expr.type isnt #Constant or node.expr.value isnt true)
			.cond(node.expr)
		.trailing(node.pos1)
		.nl()
		.print('\t{')
		.trailing(node.pos2)
		.nl()
		for (i = 0; false isnt c = node[i]; ++i)
			{
			.print("case")
			for (j = 0; j < c.size; ++j, .sep = ", ")
				.fmt(c[j])
			.print(":", nopre:)
			.trailing(c.end)
			.nl()
			.statements(c.body)
			}
		if node.def isnt false
			{
			.print("default:")
			.trailing(node.posdef)
			.nl()
			.statements(node.def)
			}
		.print("\t}")
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
			body = [body, false, type: #Compound, size: 1, pos: body.pos, end: body.end]
		.body(body)
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
	body(node)
		{
		// leading newline
		.nl()
		++.indent
		.fmt(node)
		--.indent
		// no trailing newline
		}
	Compound(node)
		{
		.print("{")
		.trailing(node.pos+1)
		.nl()

		.statements2(node)

		.leading(node, node.end - 1)
		.print("}")
		}
	statements(node)
		{
		++.indent
		.statements2(node)
		--.indent
		}
	statements2(node)
		{
		for (i = 0; false isnt stmt = node[i]; ++i)
			{
			.fmt(stmt, allowBlankLineAfter: false isnt node[i+1])
			if i+1 < node.size and node[i+1].type is #Compound and node[i+1].size is 0
				{
				.print(';;', nopre:)
				.trailing(node[i+1].pos)
				++i
				}
			.nl()
			}
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
//	println(@args)
//		{
//		.print(@args)
//		.nl()
//		}
	lineState: #start //  start, comments, other
	print(@args)
		{
		.softLineBreak = .extraLineBreak = false
		.lineState = #other
		.print_(args)
		}
	printExtra(@args)
		{
		if args[0] is '/*unused*/' and _node.type is #Param and _node.name is "unused"
			return // suppress /*unused*/ if parameter name is "unused"
		.print_(args)
		if .lineState is #start and args[0][..2] in ('//', '/*')
			.lineState = #comments
		}
	tempindent: 0
	print_(args)
		{
.trace(print: args)
		if .output[-1] is '\n' //FIXME will probably flatten
			{
			.output $= '\t'.Repeat(.tempindent + .indent + .nest)
			.sep = ""
			}
		.tempindent = 0
		if args.Extract(#nopre, false) is true
			.sep = ""
		nosuf = args.Extract(#nosuf, false)
		.prev = .output.Size()
		for arg in args
			{
			.output $= .sep $ arg
			.sep = " "
			}
		if nosuf
			.sep = ""
		}
	nl()
		{
.trace('nl')
		.output $= '\n'
		.softLineBreak = .extraLineBreak = false
		.lineState = #start
		}
	extra_nl()
		{
.trace('extra_nl')
		.blankLine()
		.tempindent = .nest is 0 ? 1 : 0
		}
	blankLine()
		{
		.output $= .sep.RemoveSuffix(' ') $ '\n'
		.sep = ""
		.lineState = #start
		}
	extraLineBreak()
		{
.trace('extraLineBreak', lineState: .lineState)
		if .lineState is #comments
			{
			.nl() // line with just comments
.trace("comment line")
			}
		else if .softLineBreak
			.extra_nl()
		else
			.extraLineBreak = true
		}
	softLineBreak()
		{
.trace('softLineBreak')
		if .extraLineBreak
			.extra_nl()
		else
			.softLineBreak = true
		}
trace(@args)
{
	args
//Print(@args)
}
	}
