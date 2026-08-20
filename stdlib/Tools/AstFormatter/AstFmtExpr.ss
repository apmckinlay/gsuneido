// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// expression and call/argument shapes for AstFormatter
AstFmtObject
	{
	ExprStmt(node, ctx, curr)
		{
		d = ctx.lastStmt is true
			? .Fmt(node.expr, ctx, curr)
			: .Expr(node.expr, ctx, curr)

		return .Style.Debug?(node.expr) ? .Root(d) : d
		}

	Expr(node, ctx, curr)
		{
		// convert PostInc/Dec to += 1 / -= 1 when result not used
		if node.type is #Unary and node.op in (#PostInc, #PostDec)
			return .Fmt(
				Object(type: #Binary, lhs: node.expr,
					op: node.op is #PostInc ? #AddEq : #SubEq,
					rhs: Object(type: #Constant, value: 1, symbol: false, pos: false,
						end: false), pos: node.pos, end: node.end), ctx, curr)
		return .Fmt(node, ctx, curr)
		}

	Constant(node, ctx, curr)
		{
		x = node.value
		if Type(x) is #AstNode
			return .Fmt(x, ctx, curr)
		if node.symbol
			return .Text('#' $ x)
		if ctx.constant is true and .Style.BareWord?(x)
			return .Text(x)
		if node.pos not in (0, false)
			return .srcConstant(node, ctx, x)
		return .bareConstant(x, ctx)
		}

	// x kept in its original source form (0xff, 1_000, the exact quoting)
	srcConstant(node, ctx, x)
		{
		src = .Src[node.pos .. node.end]
		if String?(x)
			return .StrTok(x, .Style.Quote(x, src, sym: ctx.noSym isnt true))
		if Number?(x)
			return .Tok(.Style.Num(x, src))
		return .Tok(Display(x))
		}

	bareConstant(x, ctx)
		{
		if String?(x)
			return .StrTok(x, .Style.Quote(x, Display(x), sym: ctx.noSym isnt true))
		return .Tok(Number?(x) ? .Style.Num(x, false) : Display(x))
		}

	Ident(node, ctx/*unused*/, curr/*unused*/)
		{
		return .Text(node.name)
		}

	Unary(node, ctx, curr)
		{
		ectx = .With(ctx, [constant: false]) // nested strings keep their quotes
		if node.op is #Not and node.expr.type is #In
			return .In(node.expr, ectx, curr, "not in")
		op = #(Add: '+', Sub: '-', Not: "not ", BitNot: '~', Inc: "++", PostInc: "++",
			Dec: "--", PostDec: "--", Div: "1/")
		if node.op is #LParen
			return .Cat('(', .Fmt(node.expr, ectx, curr), ')')
		if node.op in (#PostInc, #PostDec)
			return .Cat(.Fmt(node.expr, ectx, curr), op[node.op])
		return .Cat(op[node.op], .Fmt(node.expr, ectx, curr))
		}

	Binary(node, ctx, curr)
		{
		tight = node.op is #Mod and ctx.tight is true
		ectx = .With(ctx, [constant: false, tight: false])
		op = #(Eq: '=', AddEq: "+=", SubEq: "-=", CatEq: "$=", MulEq: "*=", DivEq: "/=",
			ModEq: "%=", LShiftEq: "<<=", RShiftEq: ">>=", BitOrEq: "|=", BitAndEq: "&=",
			BitXorEq: "^=", Is: is, Isnt: isnt, Match: "=~", MatchNot: "!~", Mod: '%',
			LShift: "<<", RShift: ">>", Lt: '<', Lte: "<=", Gt: '>', Gte: ">=")
		looser = node.op in (#Is, #Isnt, #Match, #MatchNot, #Lt, #Lte, #Gt, #Gte)
		l = looser ? .opnd(node.lhs, ectx, curr) : .Fmt(node.lhs, ectx, curr)
		r = looser ? .opnd(node.rhs, ectx, curr) : .Fmt(node.rhs, ectx, curr)
		return .Cat(l, tight ? op[node.op] : ' ' $ op[node.op] $ ' ', r)
		}

	opnd(e, ctx, curr, min = 0) // precedence-based spacing: AstFmtStyle Tight?/Nprec
		{
		return .Fmt(e, .With(ctx, [tight: .Style.Tight?(e, min)]), curr)
		}

	Nary(node, ctx, curr)
		{
		if node.op is #Cat and false isnt span = .Style.MultilineCatSpan(node)
			{
			.Cm.SkipTo(curr, span.end)
			return .Tok(.Src[span.pos .. span.end])
			}
		tight = ctx.tight is true
		// no #syms inside $ chains: pieces of a split string stay strings
		ectx = .With(ctx,
			node.op is #Cat
				? [constant: false, tight: false, noSym:]
				: [constant: false, tight: false])
		ops = #(And: and, Or: or, Add: '+', Cat: '$', Mul: '*', BitOr: '|', BitAnd: '&',
			BitXor: '^')
		op = ops[node.op]
		p = .Style.Nprec[node.op]
		parts = [.opnd(node[0], ectx, curr, min: p)]
		.naryTerms(parts, node, curr, [:op, :ectx, :tight, :p])
		return tight ? .Catl(parts) : .Group(.Nest(.Fill(parts)))
		}

	naryTerms(parts, node, curr, cfg)
		{
		for (i = 1; i < node.size; ++i)
			{
			st = .signedTerm(cfg.op, node[i])
			parts.Add(cfg.tight ? .Text(st.sep) : .Cat(' ' $ st.sep, .Line))
			parts.Add(st.fold
				? .Fmt(st.node, cfg.ectx, curr)
				: .opnd(st.node, cfg.ectx, curr, min: cfg.p))
			}
		}

	// fold  + -x to - x,  + <neg const> to - <const>,  * /x to / x
	signedTerm(op, e)
		{
		if op is '+'
			return .foldPlus(e)
		if op is '*' and e.type is #Unary and e.op is #Div
			return [sep: '/', node: e.expr, fold:]
		return [sep: op, node: e, fold: false]
		}

	foldPlus(e)
		{
		if e.type is #Unary and e.op is #Sub
			return [sep: '-', node: e.expr, fold:]
		if e.type is #Constant and Number?(e.value) and e.value < 0
			return [sep: '-',
				node: Object(type: #Constant, value: -e.value, symbol: false, pos: e.pos,
					end: e.end), fold:]
		return [sep: '+', node: e, fold: false]
		}

	// emit the negation of expr: swap is/isnt and =~/!~, in becomes not in,
	// double negation cancels; otherwise 'not', parenthesized unless atomic
	notExpr(expr, ctx, curr)
		{
		if expr.type is #Binary and expr.op in (#Is, #Isnt, #Match, #MatchNot)
			return .notBinary(expr, ctx, curr)
		if expr.type is #In
			return .In(expr, ctx, curr, "not in")
		if expr.type is #Unary
			return .notUnary(expr, ctx, curr)
		if expr.type in (#Ident, #Constant, #Mem, #Call)
			return .Cat("not ", .Fmt(expr, ctx, curr))
		return .Cat("not (", .Fmt(expr, ctx, curr), ')')
		}

	notBinary(expr, ctx, curr)
		{
		op = #(Is: isnt, Isnt: is, Match: "!~", MatchNot: "=~")
		l = .Fmt(expr.lhs, ctx, curr)
		r = .Fmt(expr.rhs, ctx, curr)
		return .Cat(l, ' ' $ op[expr.op] $ ' ', r)
		}

	notUnary(expr, ctx, curr)
		{
		if expr.op is #Not // double negation cancels
			return .Fmt(expr.expr, ctx, curr)
		if expr.op is #LParen
			return .Style.Negatable?(expr.expr)
				? .Cat('(', .notExpr(expr.expr, ctx, curr), ')')
				: .Cat("not ", .Fmt(expr, ctx, curr))
		return .Cat("not (", .Fmt(expr, ctx, curr), ')')
		}

	Trinary(node, ctx, curr)
		{
		ectx = .With(ctx, [constant: false])
		// cond ? true : false is just cond; ? false : true is its negation
		if .Style.BoolConst?(node.t, true) and .Style.BoolConst?(node.f, false)
			return .Fmt(node.cond, ectx, curr)
		if .Style.BoolConst?(node.t, false) and .Style.BoolConst?(node.f, true)
			return .notExpr(node.cond, ectx, curr)
		c = .Fmt(node.cond, ectx, curr)
		t = .Fmt(node.t, ectx, curr)
		f = .Fmt(node.f, ectx, curr)
		return .Group(.Cat(c, .Nest(.Cat(.Line, "? ", t, .Line, ": ", f))))
		}

	Mem(node, ctx, curr)
		{
		if node.mem.type is #Constant and String?(node.mem.value) and
			node.mem.value.Identifier?()
			return .dotMem(node, ctx, curr)
		ex = .Fmt(node.expr, ctx, curr)
		m = .opnd(node.mem, ctx, curr)
		return .Cat(ex, '[', m, ']')
		}

	dotMem(node, ctx, curr)
		{
		ectx = ctx.chain is true
			? .With(ctx, [chain: .Style.DottedCall?(node.expr)])
			: ctx
		dot = ctx.chain is true
			? .Cat('.', .Soft, node.mem.value)
			: .Text('.' $ node.mem.value)
		if .showDot?(node, ctx)
			return .Cat(.Fmt(node.expr, ectx, curr), dot)
		return .Cat(.Cm.Leading(curr, node.dotpos), '.' $ node.mem.value)
		}

	// show the receiver + dot unless it is an implicit `this.` (then just `.x`),
	// but keep an explicit this. when the method disambiguates a private member
	showDot?(node, ctx)
		{
		if node.expr.type isnt #Ident or node.expr.name isnt #this
			return true
		return ctx.method is true and node.mem.value[0].Lower?() and
			.Style.ExplicitThis?(node)
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
		ex = .Fmt(node.expr, ctx, curr)
		// drop a leading 0 in x[0::n] so it becomes x[::n] unless there is a comment that
		// hugs it, which would be orphanes otherwise and causes issues with the
		// AstFmtEquals? class' checking
		dropFrom = .Style.ZeroIdx?(node.from) and hi isnt false and
			not .Cm.CommentIn?(node.expr.end, hi.pos)
		f = dropFrom ? false : node.from
		from = f is false ? false : .opnd(f, ctx, curr)
		simple = .Style.Simple(f) and .Style.Simple(hi)
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
			parts.Add(.Fmt(arg, ctx, curr))
			}
		return .Group(.Cat(ex, ' ' $ op $ " (", .Nest(.Fill(parts)), ')'))
		}

	Call(node, ctx, curr)
		{
		parts = Object()
		open = .callOpen(node, ctx, curr, parts)
		last = node.size - 1
		blockArg = .trailingBlockArg(node, ctx, last)
		if blockArg isnt false
			--last
		docs = .argDocs(node, ctx, curr, last)
		.placeArgs(node, parts, docs, last, open.lead)
		parts.Add(.Text(open.delim))
		d = .callGroup(parts, ctx, open.chain)
		if blockArg isnt false
			d = .Cat(d,
				.FuncBody(blockArg.expr, .With(ctx, [singleLine: false, chain: false]),
					curr, extra: .BlockParams(blockArg.expr, ctx, curr)))
		return .Cat(d, .Cm.Trailing(curr, node.end, ctx.parentEnd, allowBlank:))
		}

	// emit the callee and opening delimiter; report the closing delim and the
	// layout flags (lead break, chain grouping) the rest of Call needs
	callOpen(node, ctx, curr, parts)
		{
		if false isnt newBase = .Style.NewCall(node)
			{
			// always keep the parens: new X().M() must not become new X.M()
			// which would rebind new to the whole chain
			parts.Add(
				.Cat("new ", .Fmt(newBase, .With(ctx, [chain: false]), curr), .Text('(')))
			return [delim: ')', lead:, chain: false]
			}
		if .Style.SuperNew(node)
			{
			parts.Add(.Text("super("))
			return [delim: ')', lead:, chain: false]
			}
		if .Style.UseBrackets(node)
			{
			parts.Add(.Text('['))
			return [delim: ']', lead: false, chain: false]
			}
		chain = ctx.chain is true or .Style.ChainBreak?(node)
		parts.Add(.Fmt(node.func, .With(ctx, [:chain]), curr))
		parts.Add(.Text('('))
		n = .Style.Fnlen(node.func)
		return [delim: ')', lead: n isnt false and n+1 > .TabWidth, :chain]
		}

	trailingBlockArg(node, ctx, last)
		{
		if last < 0 or ctx.cond is true
			return false
		arg = node[last]
		return arg.name is #block and arg.expr.type is #Block ? arg : false
		}

	argDocs(node, ctx, curr, last)
		{
		docs = Object()
		for (i = 0; i <= last; ++i)
			docs.Add(.argDoc(node[i], ctx, curr,
				i is last or not .Style.ArgComma?(node[i], node[i+1]) ? "" : ','))
		return docs
		}

	placeArgs(node, parts, docs, last, lead)
		{
		if not docs.Empty?()
			{
			// a block argument indents from the statement, not the arg list
			block = .blockArg?(node, last)
			if .Style.VerticalArgs?(node, last)
				parts.Add(.Nest(.Cat(.Style.LeadBreak?(node[0]) ? .Hard : false,
					.Seq(docs, .Hard))))
			else // a sole argument only earns the lead break if it then fits flat,
			 if // so Lead sits outside the arg-list nest and supplies its own
			lead and not block and docs.Size() is 1
				parts.Add(.Lead(docs[0]))
			else
				parts.Add(
					block ? .argList(node, docs) : .Nest(.argList(node, docs, :lead)))
			}
		}

	blockArg?(node, last)
		{
		for (i = 0; i <= last; ++i)
			if node[i].expr.type is #Block
				return true
		return false
		}

	callGroup(parts, ctx, chain)
		{
		if ctx.chain is true // a chain part: the chain head owns the group
			return .Catl(parts)
		return chain is true ? .Group(.Nest(.Catl(parts))) : .Group(.Catl(parts))
		}

	argList(node, docs, lead = false)
		{
		if .ArgWrap is #onePerLine
			return .Cat(.Soft, .Seq(docs, .Line))
		a = .fillArgs(node, docs)
		if lead
			{
			a.Add(.Soft, at: 0)
			a.Add(.Text(""), at: 0)
			}
		return .Fill(a)
		}

	// a break the author wrote between two arguments is kept; the rest fill
	fillArgs(node, docs)
		{
		a = Object()
		for (i = 0; i < docs.Size(); ++i)
			{
			if i > 0
				a.Add(.Style.ArgBrokeInSrc?(node, i) ? .Hard : .Line)
			a.Add(docs[i])
			}
		return a
		}

	argDoc(arg, ctx, curr, sep)
		{
		parts = Object()
		cctx = .LeadComments(arg, ctx, curr, parts)
		parts.Add(.Cat(.argument(arg, cctx, curr), sep))
		.TrailComments(arg, ctx, curr, parts)
		return .Catl(parts)
		}

	argument(arg, ctx, curr)
		{
		parts = Object()
		if .Style.Shorthand?(arg)
			parts.Add(.Text(':' $ arg.name))
		else
			.argNamed(arg, ctx, curr, parts)
		return .Catl(parts)
		}

	argNamed(arg, ctx, curr, parts)
		{
		if arg.name isnt false
			parts.Add(.Text(.argKey(arg.name)))
		// a bare `name:` (value true) emits just the key; anything else emits a value
		if arg.name is false or arg.expr.type isnt #Constant or arg.expr.value isnt true
			.argValue(arg, ctx, curr, parts)
		}

	argValue(arg, ctx, curr, parts)
		{
		if arg.name isnt false and arg.name not in ('@', "@+1")
			parts.Add(.Text(' '))
		ectx = .With(ctx,
			arg.expr.type is #Block ? [blockArg:, chain: false] : [chain: false])
		parts.Add(.Fmt(arg.expr, ectx, curr))
		}

	argKey(name)
		{
		if name in ('@', "@+1")
			return name
		if String?(name) and name.Identifier?()
			return name $ ':'
		return Display(name) $ ':'
		}
	}
