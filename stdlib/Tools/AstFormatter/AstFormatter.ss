// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// ref: Wadler, P., & Kilmer, J. (2002). A prettier printer.
AstFmtDoc
	{
	// argWrap
	//		fill = pack arguments, breaking only when needed
	//		onePerLine = every argument on each line
	// bodyStyle
	//		preserve = keep the original source code's shape
	//		fit = collapse any body to one line
	//		expand = every body breaks
	CallClass(src, width = 90, tabWidth = 4, argWrap = #fill, bodyStyle = #preserve)
		{
		Assert(Compilable?(src))
		Assert(Number?(width) and Number?(tabWidth))
		Assert(#(fill, onePerLine).Has?(argWrap))
		Assert(#(preserve, fit, expand).Has?(bodyStyle))

		// body can change layout upstream; iterate to the fixpoint
		out = (new this(width, tabWidth, argWrap, bodyStyle)).Process(src)
		for ..2
			{
			next = (new this(width, tabWidth, argWrap, bodyStyle)).Process(out)
			if next is out
				break
			out = next
			}

		// we must ensure the output code always compiles; else hard error
		Assert(Compilable?(out))
		return out
		}

	New(.width = 90, .tabWidth = 4, .argWrap = #fill, .bodyStyle = #preserve)
		{
		Assert(Number?(.width) and Number?(.tabWidth))
		Assert(#(fill, onePerLine).Has?(.argWrap))
		Assert(#(preserve, fit, expand).Has?(.bodyStyle))
		}

	maxpos: 9_999_999
	Process(src)
		{
		.src = src
		ast = Suneido.Parse(src)
		if Type(ast) isnt #AstNode // simple constant record e.g. a number
			return .rtrim(src) $ '\n'
		.cm = AstFmtComments(src)
		.style = AstFmtStyle(src, .cm)
		curr = [i: 0, done: 0, blank: false]
		ctx = [constant: false, noSym: false, cond: false, singleLine: false,
			blockArg: false, tight: false, method: false, lastStmt: false, bracket: false,
			parentEnd: .maxpos]
		doc = .Cat(.fmt(ast, ctx, curr), .cm.Leading(curr, .maxpos))
		return .rtrim(AstFmtRender(.width, .tabWidth).Render(doc)) $ '\n'
		}

	with(ctx, mods)
		{
		c = ctx.Copy()
		for m in mods.Members()
			c[m] = mods[m]
		return c
		}

	fmt(node, ctx, curr, allowBlankAfter = false)
	{
	if node is false
		return false
	unusedParam = node.type is #Param and node.name is #unused
	parts = Object()
	cctx = .leadComments(node, ctx, curr, parts, :unusedParam)
	parts.Add(this[node.type](node, cctx, curr))
	.trailComments(node, ctx, curr, parts, allowBlank: allowBlankAfter, :unusedParam)
	return .Catl(parts)
	}

	leadComments(node, ctx, curr, parts, unusedParam = false)
	{
	if curr.blank is true
		{
		parts.Add(.Blank)
		curr.blank = false
		}
	cctx = ctx
	if false isnt pos = node.pos
		{
		if node.end isnt 0
			cctx = .with(ctx, [parentEnd: node.end])
		parts.Add(.cm.Leading(curr, pos, :unusedParam))
		}
	return cctx
	}

trailComments(node, ctx, curr, parts, allowBlank = false, unusedParam = false)
	{
	if false isnt node.pos and node.end not in (0, false)
		parts.Add(.cm.Trailing(curr, node.end, ctx.parentEnd, :allowBlank, :unusedParam))
	}

	strTok(x, s) // splittable with $ when it stayed a plain quoted literal
		{
		return s[0] in ('"', "'") and .style.Plain?(x) ? .Str(s) : .Tok(s)
		}

	Object(node, ctx, curr)
		{
		delims = ctx.constant is true and
			(ctx.bracket is true or .style.BracketTable?(node))
			? "[]"
			: "()"
		return .members(node, ctx, curr, delims)
		}

	Record(node, ctx, curr)
		{
		return .members(node, ctx, curr, ctx.constant is true ? "[]" : "{}")
		}

	members(node, ctx, curr, delims)
		{
		pre = ctx.constant is true ? "" : '#'
		mctx = .with(ctx, [constant:, noSym:, bracket: false])
		vert = .style.Vertical?(node)
		pad = vert ? .style.AlignWidth(node) : false
		docs = Object()
		for (i = 0; false isnt m = node[i]; ++i)
			{
			last = false is node[i+1]
			// a key-only member keeps its comma even in vertical layout,
			// else it absorbs the next member as its value
			noComma = vert is true and (m.named isnt true or m.value isnt true)
			p = Object()
			cctx = .leadComments(m, mctx, curr, p)
			p.Add(.objectMember(m, cctx, curr, last: last or noComma, :pad))
			.trailComments(m, mctx, curr, p, allowBlank: not last)
			docs.Add(.Catl(p))
			}
		if docs.Empty?()
			return .Text(pre $ delims[0] $ delims[1])
		close = node.end in (0, false) ? false : .cm.Leading(curr, node.end - 1)
		if vert
			return .Cat(pre $ delims[0], .Nest(.Cat(.Hard, .Seq(docs, .Hard))), close,
				delims[1])
		return .Group(
			.Cat(pre $ delims[0], .Nest(.Fillsep(docs, .Line)), close, delims[1]))
		}

	objectMember(m, ctx, curr, last = false, pad = false)
		{
		parts = Object()
		if not m.named
			parts.Add(.constant(m.value, ctx, curr, s: .style.MemberSrc(m), split:))
		else
			{
			parts.Add(.constant(m.key, ctx, curr))
			if m.value is true
				parts.Add(.Text(':'))
			else
				{
				parts.Add(
					pad is false
						? .Text(": ")
						: .Text(':' $ ' '.Repeat(pad - m.key.Size() + 1)))
				parts.Add(
					.constant(m.value, ctx, curr, s: .style.MemberSrc(m, keyed:), split:))
				}
			}
		if not last
			parts.Add(.Text(','))
		return .Catl(parts)
		}

	constant(x, ctx, curr, s = false, split = false)
		{
		if Type(x) is #AstNode
			return .fmt(x, s isnt false and s[0] is '[' ? .with(ctx, [bracket:]) : ctx,
				curr)
		if ctx.constant is true and .style.BareWord?(x)
			return .Text(x)
		if String?(x)
			{
			q = .style.Quote(x, s is false ? Display(x) : s, sym: ctx.noSym isnt true)
			return split is true ? .strTok(x, q) : .Tok(q)
			}
		if Number?(x) and s isnt false // Display would decimal-ize 0xff
			return .Tok(s)
		return .Tok(Display(x))
		}

	Class(node, ctx, curr)
		{
		base = .Text(node.base)
		t1 = .cm.Trailing(curr, node.pos1)
		mctx = .with(ctx, [constant:, noSym:, method:])
		parts = [.Text('{'), .cm.Trailing(curr, node.pos2)]
		for (i = 0; false isnt m = node[i]; ++i)
			{
			parts.Add(.Hard)
			last = false is node[i+1]
			p = Object()
			cctx = .leadComments(m, mctx, curr, p)
			p.Add(.classMember(m, cctx, curr))
			.trailComments(m, mctx, curr, p, allowBlank: not last)
			parts.Add(.Catl(p))
			// a blank line after every method; curr.blank dedups with any in source
			if not last and Type(m.value) is #AstNode and m.value.type is #Function
				curr.blank = true
			}
		parts.Add(.Hard)
		parts.Add(.cm.Leading(curr, node.end - 1))
		parts.Add(.Text('}'))
		return .Cat(base, t1, .Nest(.Cat(.Hard, .Catl(parts))))
		}

	classMember(m, ctx, curr)
		{
		if Type(m.value) is #AstNode
			{
			if m.value.type is #Function
				return .Cat(.Text(m.key), .func(m.value, ctx, curr))
			k = .constant(m.key, ctx, curr)
			s = .style.MemberSrc(m, keyed:)
			v = .fmt(m.value, s isnt false and s[0] is '[' ? .with(ctx, [bracket:]) : ctx,
				curr)
			return .Cat(k, ": ", v)
			}
		if String?(m.value)
			{
			s = .style.MemberSrc(m, keyed:)
			v = .style.Quote(m.value, s is false ? Display(m.value) : s)
			return .Cat(.constant(m.key, ctx, curr), ": ", .strTok(m.value, v))
			}
		return .Cat(.constant(m.key, ctx, curr), ": ", .Tok(Display(m.value)))
		}

	Function(node, ctx, curr)
		{
		return .Cat("function", .func(node, ctx, curr))
		}

	func(node, ctx, curr)
		{
		p = .params(node.params, ctx, curr, "()")
		t1 = .cm.Trailing(curr, node.pos1)
		return .Cat(p, t1, .funcBody(node, ctx, curr, pos2: node.pos2, lead: .Line))
		}

	Block(node, ctx, curr)
		{
		bctx = .with(ctx, [blockArg: false])
		if ctx.blockArg is true
			return .funcBody(node, bctx, curr, extra: .blockParams(node, bctx, curr),
				lead: .Soft)
		extra = .blockParams(node, bctx, curr)
		single = .single?(node, ctx)
		sctx = .with(bctx, [constant: false, noSym: false, singleLine:])
		br = single ? .Line : .Hard
		sep = single ? .Semi : .Hard
		stmts = .stmtDocs(node, sctx, curr, allowBlank: not single)
		inner = Object()
		for (j = 0; j < stmts.Size(); ++j)
			{
			inner.Add(j is 0 ? br : sep)
			inner.Add(stmts[j])
			}
		inner.Add(br)
		inner.Add(.cm.Leading(curr, node.end))
		inner.Add(.Text('}'))
		d = .Cat(.Text('{'), extra, .Nest(.Catl(inner)))
		return single ? .Group(d) : d
		}

	blockParams(node, ctx, curr)
		{
		params = node.params
		return params.size > 1 or (params.size is 1 and params[0].name isnt #it)
			? .params(params, ctx, curr, "||")
			: false
		}

	params(params, ctx, curr, delims)
		{
		pctx = .with(ctx, [constant: false, noSym: false, singleLine:])
		docs = Object()
		for (i = 0; false isnt p = params[i]; ++i)
			docs.Add(.fmt(p, pctx, curr))
		if docs.Empty?()
			return .Text(delims[0] $ delims[1])
		return .Group(.Cat(delims[0], .Nest(.Fillsep(docs, .Cat(',', .Line))), delims[1]))
		}

	Param(node, ctx, curr)
		{
		parts = [.Text(node.name)]
		if node.hasdef
			{
			parts.Add(
				.cm.Trailing(curr, node.pos + node.name.Size(),
					unusedParam: node.name is #unused))
			parts.Add(.Text(" = "))
			s = node.pos in (0, false) or node.end in (0, false)
				? false
				: .style.ValueSrc(.src[node.pos .. node.end].AfterFirst('='))
			parts.Add(.constant(node.defval, ctx, curr, :s))
			}
		return .Catl(parts)
		}

	funcBody(node, ctx, curr, extra = false, pos2 = false, lead = false)
		{
		single = .single?(node, ctx)
		bctx = .with(ctx, [constant: false, noSym: false, singleLine:])
		br = single ? .Line : .Hard
		sep = single ? .Semi : .Hard
		parts = [single ? (lead is false ? .Soft : lead) : .Hard, .Text('{')]
		if extra isnt false
			parts.Add(extra)
		if pos2 isnt false
			parts.Add(.cm.Trailing(curr, pos2))
		stmts = .stmtDocs(node, bctx, curr, allowBlank: not single)
		for (j = 0; j < stmts.Size(); ++j)
			{
			parts.Add(j is 0 ? br : sep)
			parts.Add(stmts[j])
			}
		parts.Add(br)
		parts.Add(.cm.Leading(curr, node.end))
		parts.Add(.Text('}'))
		d = .Nest(.Catl(parts))
		return single ? .Group(d) : d
		}

	single?(node, ctx)
		{
		if ctx.singleLine isnt true or .bodyStyle is #expand
			return false
		return .bodyStyle is #fit
			? .style.OkToSingleLine2(node)
			: .style.OkToSingleLine(node)
		}

	ExprStmt(node, ctx, curr)
		{
		d = ctx.lastStmt is true ? .fmt(node.expr, ctx, curr) : .expr(node.expr, ctx, curr)

		return .style.Debug?(node.expr) ? .Root(d) : d
		}

	expr(node, ctx, curr)
		{
		// convert PostInc/Dec to pre Inc/Dec when result not used
		if node.type is #Unary and node.op in (#PostInc, #PostDec)
			return .fmt(
				Object(type: #Unary, expr: node.expr, op: node.op.RemovePrefix(#Post),
					pos: node.pos, end: node.end), ctx, curr)
		return .fmt(node, ctx, curr)
		}

	Constant(node, ctx, curr)
		{
		x = node.value
		if Type(x) is #AstNode
			return .fmt(x, ctx, curr)
		if node.symbol
			return .Text('#' $ x)
		if ctx.constant is true and .style.BareWord?(x)
			return .Text(x)
		if node.pos not in (0, false) and String?(x)
			return .strTok(x,
				.style.Quote(x, .src[node.pos .. node.end], sym: ctx.noSym isnt true))
		if node.pos not in (0, false) and Number?(x)
			return .Tok(.src[node.pos .. node.end])
		if String?(x)
			return .strTok(x, .style.Quote(x, Display(x), sym: ctx.noSym isnt true))
		return .Tok(Display(x))
		}

	Ident(node, ctx/*unused*/, curr/*unused*/)
		{
		return .Text(node.name)
		}

	Unary(node, ctx, curr)
		{
		ectx = .with(ctx, [constant: false]) // nested strings keep their quotes
		if node.op is #Not and node.expr.type is #In
			return .In(node.expr, ectx, curr, "not in")
		op = #(Add: '+', Sub: '-', Not: "not ", BitNot: '~', Inc: "++", PostInc: "++",
			Dec: "--", PostDec: "--", Div: "1/")
		if node.op is #LParen
			return .Cat('(', .fmt(node.expr, ectx, curr), ')')
		if node.op in (#PostInc, #PostDec)
			return .Cat(.fmt(node.expr, ectx, curr), op[node.op])
		return .Cat(op[node.op], .fmt(node.expr, ectx, curr))
		}

	Binary(node, ctx, curr)
		{
		tight = node.op is #Mod and ctx.tight is true
		ectx = .with(ctx, [constant: false, tight: false])
		op = #(Eq: '=', AddEq: "+=", SubEq: "-=", CatEq: "$=", MulEq: "*=", DivEq: "/=",
			ModEq: "%=", LShiftEq: "<<=", RShiftEq: ">>=", BitOrEq: "|=", BitAndEq: "&=",
			BitXorEq: "^=", Is: is, Isnt: isnt, Match: "=~", MatchNot: "!~", Mod: '%',
			LShift: "<<", RShift: ">>", Lt: '<', Lte: "<=", Gt: '>', Gte: ">=")
		looser = node.op in (#Is, #Isnt, #Match, #MatchNot, #Lt, #Lte, #Gt, #Gte)
		l = looser ? .opnd(node.lhs, ectx, curr) : .fmt(node.lhs, ectx, curr)
		r = looser ? .opnd(node.rhs, ectx, curr) : .fmt(node.rhs, ectx, curr)
		return .Cat(l, tight ? op[node.op] : ' ' $ op[node.op] $ ' ', r)
		}

	opnd(e, ctx, curr, min = 0) // precedence-based spacing: AstFmtStyle Tight?/Nprec
		{
		return .fmt(e, .with(ctx, [tight: .style.Tight?(e, min)]), curr)
		}

	Nary(node, ctx, curr)
		{
		tight = ctx.tight is true
		// no #syms inside $ chains: pieces of a split string stay strings
		ectx = .with(ctx,
			node.op is #Cat
				? [constant: false, tight: false, noSym:]
				: [constant: false, tight: false])
		ops = #(And: and, Or: or, Add: '+', Cat: '$', Mul: '*', BitOr: '|', BitAnd: '&',
			BitXor: '^')
		op = ops[node.op]
		p = .style.Nprec[node.op]
		parts = [.opnd(node[0], ectx, curr, min: p)]
		for (i = 1; i < node.size; ++i)
			{
			st = .signedTerm(op, node[i])
			parts.Add(tight ? .Text(st.sep) : .Cat(' ' $ st.sep, .Line))
			parts.Add(
				st.fold ? .fmt(st.node, ectx, curr) : .opnd(st.node, ectx, curr, min: p))
			}
		return tight ? .Catl(parts) : .Group(.Nest(.Fill(parts)))
		}

	// fold  + -x to - x,  + <neg const> to - <const>,  * /x to / x
	signedTerm(op, e)
		{
		if op is '+' and e.type is #Unary and e.op is #Sub
			return [sep: '-', node: e.expr, fold:]
		if op is '+' and e.type is #Constant and Number?(e.value) and e.value < 0
			return [sep: '-',
				node: Object(type: #Constant, value: -e.value, symbol: false, pos: e.pos,
					end: e.end), fold:]
		if op is '*' and e.type is #Unary and e.op is #Div
			return [sep: '/', node: e.expr, fold:]
		return [sep: op, node: e, fold: false]
		}

	// emit the negation of expr: swap is/isnt and =~/!~, in becomes not in,
	// double negation cancels; otherwise 'not', parenthesized unless atomic
	notExpr(expr, ctx, curr)
		{
		if expr.type is #Binary and expr.op in (#Is, #Isnt, #Match, #MatchNot)
			{
			op = #(Is: isnt, Isnt: is, Match: "!~", MatchNot: "=~")
			l = .fmt(expr.lhs, ctx, curr)
			r = .fmt(expr.rhs, ctx, curr)
			return .Cat(l, ' ' $ op[expr.op] $ ' ', r)
			}
		if expr.type is #In
			return .In(expr, ctx, curr, "not in")
		if expr.type is #Unary and expr.op is #Not
			return .fmt(expr.expr, ctx, curr)
		if expr.type is #Unary and expr.op is #LParen
			return .style.Negatable?(expr.expr)
				? .Cat('(', .notExpr(expr.expr, ctx, curr), ')')
				: .Cat("not ", .fmt(expr, ctx, curr))
		if expr.type in (#Ident, #Constant, #Mem, #Call)
			return .Cat("not ", .fmt(expr, ctx, curr))
		return .Cat("not (", .fmt(expr, ctx, curr), ')')
		}

	Trinary(node, ctx, curr)
		{
		ectx = .with(ctx, [constant: false])
		// cond ? true : false is just cond; ? false : true is its negation
		if .style.BoolConst?(node.t, true) and .style.BoolConst?(node.f, false)
			return .fmt(node.cond, ectx, curr)
		if .style.BoolConst?(node.t, false) and .style.BoolConst?(node.f, true)
			return .notExpr(node.cond, ectx, curr)
		c = .fmt(node.cond, ectx, curr)
		t = .fmt(node.t, ectx, curr)
		f = .fmt(node.f, ectx, curr)
		return .Group(.Cat(c, .Nest(.Cat(.Line, "? ", t, .Line, ": ", f))))
		}

	Mem(node, ctx, curr)
		{
		if node.mem.type is #Constant and String?(node.mem.value) and
			node.mem.value.Identifier?()
			{
			if node.expr.type isnt #Ident or node.expr.name isnt #this
				return .Cat(.fmt(node.expr, ctx, curr), '.' $ node.mem.value)
			if ctx.method is true and node.mem.value[0].Lower?() and
				.style.ExplicitThis?(node.expr)
				return .Cat(.fmt(node.expr, ctx, curr), '.' $ node.mem.value)
			return .Cat(.cm.Leading(curr, node.dotpos), '.' $ node.mem.value)
			}
		ex = .fmt(node.expr, ctx, curr)
		m = .opnd(node.mem, ctx, curr)
		return .Cat(ex, '[', m, ']')
		}

	RangeTo(node, ctx, curr)
		{
		return .range(node, ctx, curr, node.to, "..")
		}

	RangeLen(node, ctx, curr)
		{
		return .range(node, ctx, curr, node.len, "::")
		}

	range(node, ctx, curr, hi, dots) // x[from .. to] and x[from :: len] share a shape
		{
		ex = .fmt(node.expr, ctx, curr)
		f = .style.ZeroIdx?(node.from) and hi isnt false ? false : node.from
		from = f is false ? false : .opnd(f, ctx, curr)
		simple = .style.Simple(f) and .style.Simple(hi)
		to = hi is false ? false : .opnd(hi, ctx, curr)
		sep = simple ? dots : .rangeSep(from, to, dots)
		return .Cat(ex, '[', from, sep, to, ']')
		}

	rangeSep(from, to, dots)
		{
		return (from is false ? "" : ' ') $ dots $ (to is false ? "" : ' ')
		}

	In(node, ctx, curr, op = #in)
		{
		ex = .opnd(node.expr, ctx, curr)
		parts = Object()
		for (i = 0; false isnt arg = node[i]; ++i)
			{
			if i > 0
				parts.Add(.Cat(',', .Line))
			parts.Add(.fmt(arg, ctx, curr))
			}
		return .Group(.Cat(ex, ' ' $ op $ " (", .Nest(.Fill(parts)), ')'))
		}

	Call(node, ctx, curr)
		{
		parts = Object()
		delim = ')'
		lead = false
		if .style.NewThis(node)
			{
			if node.size is 0
				return .Cat("new this",
					.cm.Trailing(curr, node.end, ctx.parentEnd, allowBlank:))
			parts.Add(.Text("new this("))
			lead = true
			}
		else if .style.SuperNew(node)
			{
			parts.Add(.Text("super("))
			lead = true
			}
		else if .style.UseBrackets(node)
			{
			parts.Add(.Text('['))
			delim = ']'
			}
		else
			{
			parts.Add(.fmt(node.func, ctx, curr))
			parts.Add(.Text('('))
			n = .style.Fnlen(node.func)
			lead = n isnt false and n+1 > .tabWidth
			}
		last = node.size - 1
		blockArg = false
		if last >= 0 and ctx.cond isnt true
			{
			arg = node[last]
			if arg.name is #block and arg.expr.type is #Block
				{
				blockArg = arg
				--last
				}
			}
		docs = Object()
		for (i = 0; i <= last; ++i)
			docs.Add(
				.argDoc(node[i], ctx, curr,
					i is last or not .style.ArgComma?(node[i], node[i+1]) ? "" : ','))
		if not docs.Empty?()
			{
			// a lone block argument indents from the statement, not the arg list
			loneBlock = node.size is 1 and node[0].name is false and
				node[0].expr.type is #Block
			parts.Add(loneBlock ? .argList(docs) : .Nest(.argList(docs, :lead)))
			}
		parts.Add(.Text(delim))
		d = .Group(.Catl(parts))
		if blockArg isnt false
			d = .Cat(d,
				.funcBody(blockArg.expr, .with(ctx, [singleLine: false]), curr,
					extra: .blockParams(blockArg.expr, ctx, curr)))
		return .Cat(d, .cm.Trailing(curr, node.end, ctx.parentEnd, allowBlank:))
		}

	argList(docs, lead = false)
		{
		if .argWrap is #onePerLine
			return .Cat(.Soft, .Seq(docs, .Line))
		a = .Interleave(docs, .Line)
		if lead
			{
			a.Add(.Soft, at: 0)
			a.Add(.Text(""), at: 0)
			}
		return .Fill(a)
		}

	argDoc(arg, ctx, curr, sep)
		{
		parts = Object()
		cctx = .leadComments(arg, ctx, curr, parts)
		parts.Add(.Cat(.argument(arg, cctx, curr), sep))
		.trailComments(arg, ctx, curr, parts)
		return .Catl(parts)
		}

	argument(arg, ctx, curr)
		{
		parts = Object()
		if .style.Shorthand?(arg)
			parts.Add(.Text(':' $ arg.name))
		else
			{
			if arg.name isnt false
				parts.Add(.Text(.argKey(arg.name)))
			if arg.name is false or arg.expr.type isnt #Constant or
				arg.expr.value isnt true
				{
				if arg.name isnt false and arg.name not in ('@', "@+1")
					parts.Add(.Text(' '))
				ectx = arg.expr.type is #Block ? .with(ctx, [blockArg:]) : ctx
				parts.Add(.fmt(arg.expr, ectx, curr))
				}
			}
		return .Catl(parts)
		}

	argKey(name)
		{
		if name in ('@', "@+1")
			return name
		if String?(name) and name.Identifier?()
			return name $ ':'
		return Display(name) $ ':'
		}

	Return(node, ctx, curr)
		{
		docs = Object()

		for (i = 0; false isnt e = node[i]; ++i)
			docs.Add(.fmt(e, ctx, curr))

		kw = node.throw is true ? "return throw" : #return

		if docs.Empty?()
			return .Text(kw)

		return .Cat(kw $ ' ', .Seq(docs, .Text(", ")))
		}

	Throw(node, ctx, curr)
		{
		return .Cat("throw ", .fmt(node.expr, ctx, curr))
		}

	TryCatch(node, ctx, curr)
		{
		parts = [.Text(#try)]
		parts.Add(.cm.Trailing(curr, node.pos + 3))/*= "try".Size() */
		parts.Add(
			.ctlBody(node.try, ctx, curr, guard: node.catch is false ? false : #Catch))
		if false isnt c = node.catch
			{
			parts.Add(.Hard)
			parts.Add(.Text(#catch))
			if false isnt var = .style.CatchVar(node)
				{
				v = .fmt(var, ctx, curr)
				pat = node.catchpat is false ? "" : ", " $ Display(node.catchpat)
				parts.Add(.Cat(" (", v, pat, ')'))
				}
			parts.Add(.cm.Trailing(curr, node.catchend))
			parts.Add(.ctlBody(c, ctx, curr))
			}
		return .Catl(parts)
		}

	If(node, ctx, curr)
		{
		parts = Object()
		end = false
		forever
			{
			c = .cond(node.cond, ctx, curr)
			parts.Add(
				.Cat("if ", c,
					.ctlBody(node.t, ctx, curr, guard: node.f is false ? false : #Else)))
			if node.f is false
				break
			parts.Add(.Hard)
			parts.Add(.Text(#else))
			parts.Add(.cm.Trailing(curr, node.elseend))
			f = node.f
			if false isnt inner = .style.ElseChain(f, curr)
				{
				if end is false or f.end > end
					end = f.end
				f = inner
				}
			if f.type is #If
				{
				parts.Add(.Text(' '))
				node = f
				continue
				}
			parts.Add(.ctlBody(f, ctx, curr))
			break
			}
		if end isnt false
			.cm.SkipTo(curr, end)
		return .Catl(parts)
		}

	cond(expr, ctx, curr)
		{
		return .fmt(expr, .with(ctx, [cond:]), curr)
		}

	ctlBody(node, ctx, curr, guard = false, brace = false)
		{
		bctx = .with(ctx, [lastStmt: false])
		if brace is false and false isnt stmt = .style.Unbrace(node, guard)
			{
			d = .fmt(stmt, bctx, curr)
			.cm.SkipTo(curr, node.end)
			return .Nest(.Cat(.Hard, d))
			}
		return .Nest(.Cat(.Hard, .fmt(node, bctx, curr)))
		}

	Switch(node, ctx, curr)
		{
		parts = [.Text(#switch)]
		e = node.expr is false ? false : .style.SwitchExpr(node.expr)
		if e isnt false and (e.type isnt #Constant or e.value isnt true)
			{
			parts.Add(.Text(' '))
			parts.Add(.cond(e, ctx, curr))
			}
		parts.Add(.cm.Trailing(curr, node.pos1))
		parts.Add(.Hard)
		parts.Add(.Text('\t{'))
		parts.Add(.cm.Trailing(curr, node.pos2))
		for (i = 0; false isnt c = node[i]; ++i)
			{
			parts.Add(.Hard)
			vals = Object()
			for (j = 0; j < c.size; ++j)
				{
				if j > 0
					vals.Add(.Text(", "))
				vals.Add(.fmt(c[j], ctx, curr))
				}
			parts.Add(.Cat("case ", .Catl(vals), ':'))
			parts.Add(.cm.Trailing(curr, c.end))
			.caseBody(parts, c.body, ctx, curr)
			}
		if node.def isnt false
			{
			parts.Add(.Hard)
			parts.Add(.Text("default:"))
			parts.Add(.cm.Trailing(curr, node.posdef))
			.caseBody(parts, node.def, ctx, curr)
			}
		parts.Add(.Hard)
		parts.Add(.Text('\t}'))
		return .Catl(parts)
		}

	caseBody(parts, body, ctx, curr)
		{
		stmts = .stmtDocs(body, ctx, curr)
		if not stmts.Empty?()
			parts.Add(.Nest(.Cat(.Hard, .Seq(stmts, .Hard))))
		}

	Forever(node, ctx, curr)
		{
		return .Cat(#forever, .ctlBody(node.body, ctx, curr))
		}

	While(node, ctx, curr)
		{
		if node.cond.type is #Constant and node.cond.value is true
			return .Forever(node, ctx, curr)
		c = .cond(node.cond, ctx, curr)
		return .Cat("while ", c, .ctlBody(node.body, ctx, curr))
		}

	DoWhile(node, ctx, curr)
		{
		body = node.body
		if body.type isnt #Compound
			body = [body, false, type: #Compound, size: 1, pos: body.pos, end: body.end]
		b = .ctlBody(body, ctx, curr, brace:)
		c = .cond(node.cond, ctx, curr)
		return .Cat(#do, b, " while ", c)
		}

	ForIn(node, ctx, curr)
		{
		// range form (for ..n, for i in a..b) has children [expr, expr2, body];
		// expr2 must not be accessed when absent, so emit the header as is
		if false isnt node.children[2]
			{
			header = .rtrim(.src[node.pos .. node.body.pos])
			.cm.SkipTo(curr, node.body.pos)
			return .Cat(.Tok(header), .ctlBody(node.body, ctx, curr))
			}
		v = node.var
		if node.var2 isnt ""
			v $= ", " $ node.var2
		ex = .cond(node.expr, ctx, curr)
		return .Cat("for " $ v $ " in ", ex, .ctlBody(node.body, ctx, curr))
		}

	rtrim(s)
		{
		while s.Size() > 0 and s[-1] in (' ', '\t', '\n', '\r')
			s = s[..-1]
		return s
		}

	For(node, ctx, curr)
		{
		if node.init.Empty?() and node.cond is false and node.inc.Empty?()
			return .Forever(node, ctx, curr)
		parts = [.Text("for (")]
		sep = ""
		for e in node.init
			{
			parts.Add(.Text(sep))
			parts.Add(.expr(e, ctx, curr))
			sep = ", "
			}
		parts.Add(.Text(';'))
		if node.cond isnt false
			{
			parts.Add(.Text(' '))
			parts.Add(.cond(node.cond, ctx, curr))
			}
		parts.Add(.Text(';'))
		sep = ' '
		for e in node.inc
			{
			parts.Add(.Text(sep))
			parts.Add(.expr(e, ctx, curr))
			sep = ", "
			}
		parts.Add(.Text(')'))
		return .Cat(.Catl(parts), .ctlBody(node.body, ctx, curr))
		}

	Compound(node, ctx, curr)
		{
		open = .cm.Trailing(curr, node.pos + 1)
		stmts = .stmtDocs(node, ctx, curr)
		parts = [.Text('{'), open]
		if not stmts.Empty?()
			{
			parts.Add(.Hard)
			parts.Add(.Seq(stmts, .Hard))
			}
		parts.Add(.Hard)
		parts.Add(.cm.Leading(curr, node.end - 1))
		parts.Add(.Text('}'))
		return .Catl(parts)
		}

	stmtDocs(node, ctx, curr, allowBlank = true)
		{
		docs = Object()
		for (i = 0; false isnt stmt = node[i]; ++i)
			{
			last = false is node[i+1]
			d = .fmt(stmt, .with(ctx, [lastStmt: last]), curr,
				allowBlankAfter: allowBlank and not last)
			if i+1 < node.size and node[i+1].type is #Compound and node[i+1].size is 0
				{
				d = .Cat(d, ";;", .cm.Trailing(curr, node[i+1].pos))
				++i
				}
			docs.Add(d)
			}
		return docs
		}

	Break(unused, ctx/*unused*/, curr/*unused*/)
		{
		return .Text(#break)
		}

	Continue(unused, ctx/*unused*/, curr/*unused*/)
		{
		return .Text(#continue)
		}

	Default(type, node, ctx/*unused*/, curr)
		{
		if Type(node) is #AstNode and node.pos isnt false and node.end not in (0, false)
			return .verbatimNode(node, curr)
		throw "AstFormatter: cannot format " $ Display(type)
		}

	verbatimNode(node, curr)
		{
		.cm.SkipTo(curr, node.end)
		return .Tok(.src[node.pos .. node.end])
		}
	}
