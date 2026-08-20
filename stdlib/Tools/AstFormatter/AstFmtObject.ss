// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// object, record, class and function/param shapes for AstFormatter
AstFmtCore
	{
	Object(node, ctx, curr)
		{
		delims = ctx.constant is true and
			(ctx.bracket is true or .Style.BracketTable?(node))
			? "[]"
			: "()"
		return .members(node, ctx, curr, delims)
		}

	Record(node, ctx, curr)
		{
		// [..] only stays a Record when all members are named; a positional
		// member would reparse [pos..] as an Object, so keep the {..} form
		brackets = ctx.constant is true and .Style.RecordBrackets?(node)
		return .members(node, ctx, curr, brackets ? "[]" : "{}")
		}

	members(node, ctx, curr, delims)
		{
		pre = ctx.constant is true ? "" : '#'
		if ctx.cols isnt false
			return .tableRow(node, ctx, curr, pre $ delims[0], delims[1])
		vert = .Style.Vertical?(node)
		pad = vert ? .Style.AlignWidth(node) : false
		cols = vert ? .Style.TableColumns(node, .Width, .TabWidth) : false
		mctx = .With(ctx, [constant:, noSym:, bracket: false, chain: false, :cols])
		docs = .memberDocs(node, mctx, curr, pad)
		if docs.Empty?()
			return .Text(pre $ delims[0] $ delims[1])
		close = node.end in (0, false) ? false : .Cm.Leading(curr, node.end - 1)
		if vert
			// the close delimiter gets its own line, indented with the members
			return .Cat(pre $ delims[0],
				.Nest(.Cat(.Hard, .Seq(docs, .Hard), .Hard, close)), delims[1])
		return .Group(
			.Cat(pre $ delims[0], .Nest(.fillMembers(node, docs)), close, delims[1]))
		}

	memberDocs(node, mctx, curr, pad)
		{
		docs = Object()
		for (i = 0; false isnt m = node[i]; ++i)
			{
			last = false is node[i+1]
			p = Object()
			cctx = .LeadComments(m, mctx, curr, p)
			p.Add(.objectMember(m, cctx, curr, :last, :pad))
			.TrailComments(m, mctx, curr, p, allowBlank: not last)
			docs.Add(.Catl(p))
			}
		return docs
		}

	fillMembers(node, docs)
		{
		a = Object()
		for (i = 0; i < docs.Size(); ++i)
			{
			if i > 0
				a.Add(.brokeInSrc?(node, i) ? .Hard : .Line)
			a.Add(docs[i])
			}
		return .Fill(a)
		}

	brokeInSrc?(node, i)
		{
		prev = node[i-1].end
		pos = node[i].pos
		if prev in (0, false) or pos in (0, false)
			return false
		return .Src[prev..pos].Has?('\n')
		}

	objectMember(m, ctx, curr, last = false, pad = false)
		{
		parts = Object()
		src = m.named ? .Style.MemberSrc(m, keyed:) : .Style.MemberSrc(m)
		.memberValue(m, ctx, curr, src, pad, parts)
		if not last
			parts.Add(.Text(','))
		return .Catl(parts)
		}

	memberValue(m, ctx, curr, src, pad, parts)
		{
		if Type(m.value) is #AstNode and src isnt false and src.Has?('`')
			{
			if m.named
				parts.Add(.Cat(.constant(m.key, ctx, curr), ": "))
			.Cm.SkipTo(curr, m.end)
			parts.Add(.Tok(src))
			}
		else if not m.named
			parts.Add(.constant(m.value, ctx, curr, s: .Style.MemberSrc(m), split:))
		else
			.namedMember(m, ctx, curr, pad, parts)
		}

	namedMember(m, ctx, curr, pad, parts)
		{
		parts.Add(.constant(m.key, ctx, curr))
		if m.value is true
			parts.Add(.Text(':'))
		else
			{
			parts.Add(pad is false
				? .Text(": ")
				: .Text(':' $ ' '.Repeat(pad - m.key.Size() + 1)))
			parts.Add(
				.constant(m.value, ctx, curr, s: .Style.MemberSrc(m, keyed:), split:))
			}
		}

	// one row of an aligned table: unbreakable, cells padded to the column widths
	tableRow(node, ctx, curr, open, close)
		{
		rctx = .With(ctx, [cols: false])
		parts = [.Text(open)]
		for (i = 0; false isnt m = node[i]; ++i)
			{
			p = Object()
			.LeadComments(m, rctx, curr, p)
			cell = .Style.CellText(m)
			if false isnt node[i+1]
				cell $= ',' $ ' '.Repeat(ctx.cols[i] - cell.Size() + 1)
			p.Add(.Tok(cell))
			.TrailComments(m, rctx, curr, p)
			parts.Add(.Catl(p))
			}
		if node.end not in (0, false)
			parts.Add(.Cm.Leading(curr, node.end - 1))
		parts.Add(.Text(close))
		return .Catl(parts)
		}

	constant(x, ctx, curr, s = false, split = false)
		{
		if Type(x) is #AstNode
			return .fmtBracketed(x, ctx, curr, s)
		if ctx.constant is true and .Style.BareWord?(x)
			return .Text(x)
		if String?(x)
			return .strConstant(x, ctx, s, split)
		if Number?(x) // source form kept: Display would decimal-ize 0xff
			return .Tok(.Style.Num(x, s))
		return .Tok(Display(x))
		}

	// a bracketed source ('[...]') keeps the bracket context so it round-trips
	fmtBracketed(x, ctx, curr, s)
		{
		return .Fmt(x, s isnt false and s[0] is '[' ? .With(ctx, [bracket:]) : ctx, curr)
		}

	strConstant(x, ctx, s, split)
		{
		q = .Style.Quote(x, s is false ? Display(x) : s, sym: ctx.noSym isnt true)
		return split is true ? .StrTok(x, q) : .Tok(q)
		}

	StrTok(x, s) // splittable with $ when it stayed a plain quoted literal
		{
		return s[0] in ('"', "'") and .Style.Plain?(x) ? .Str(s) : .Tok(s)
		}

	Class(node, ctx, curr)
		{
		base = .Text(node.base)
		t1 = .Cm.Trailing(curr, node.pos1)
		mctx = .With(ctx, [constant:, noSym:, method:])
		pads = .Style.ClassAlign(node)
		parts = [.Text('{'), .Cm.Trailing(curr, node.pos2)]
		for (i = 0; false isnt m = node[i]; ++i)
			{
			parts.Add(.Hard)
			last = false is node[i+1]
			p = Object()
			cctx = .LeadComments(m, mctx, curr, p)
			p.Add(.classMember(m, cctx, curr, pads[i]))
			.TrailComments(m, mctx, curr, p, allowBlank: not last)
			parts.Add(.Catl(p))
			// a blank line after every method; curr.blank dedups with any in source
			if not last and Type(m.value) is #AstNode and m.value.type is #Function
				curr.blank = true
			}
		parts.Add(.Hard)
		parts.Add(.Cm.Leading(curr, node.end - 1))
		parts.Add(.Text('}'))
		return .Cat(base, t1, .Nest(.Cat(.Hard, .Catl(parts))))
		}

	classMember(m, ctx, curr, pad = false)
		{
		if Type(m.value) is #AstNode
			return .classMemberNode(m, ctx, curr)
		// align the value column across a run of scalar members
		sep = pad is false ? ": " : ':' $ ' '.Repeat(pad - m.key.Size() + 1)
		return .classMemberScalar(m, ctx, curr, sep)
		}

	classMemberNode(m, ctx, curr)
		{
		if m.value.type is #Function
			return .Cat(.Text(m.key), .func(m.value, ctx, curr))
		k = .constant(m.key, ctx, curr)
		// a constant member value has no sub-positions, so a backquoted
		// piece would be re-quoted (and its layout lost); backquoted text
		// must never be touched, so emit the whole value verbatim
		s = .Style.MemberSrc(m, keyed:)
		if s isnt false and s.Has?('`')
			{
			.Cm.SkipTo(curr, m.end)
			return .Cat(.constant(m.key, ctx, curr), ": ", .Tok(s))
			}
		return .Cat(k, ": ", .fmtBracketed(m.value, ctx, curr, s))
		}

	classMemberScalar(m, ctx, curr, sep)
		{
		if String?(m.value)
			{
			s = .Style.MemberSrc(m, keyed:)
			v = .Style.Quote(m.value, s is false ? Display(m.value) : s)
			return .Cat(.constant(m.key, ctx, curr), sep, .StrTok(m.value, v))
			}
		if Number?(m.value) // source form kept: 0xff, 1_000
			return .Cat(.constant(m.key, ctx, curr), sep,
				.Tok(.Style.Num(m.value, .Style.MemberSrc(m, keyed:))))
		return .Cat(.constant(m.key, ctx, curr), sep, .Tok(Display(m.value)))
		}

	Function(node, ctx, curr)
		{
		return .Cat("function", .func(node, ctx, curr))
		}

	func(node, ctx, curr)
		{
		p = .Params(node.params, ctx, curr, "()")
		t1 = .Cm.Trailing(curr, node.pos1)
		return .Cat(p, t1, .FuncBody(node, ctx, curr, pos2: node.pos2, lead: .Line))
		}

	Block(node, ctx, curr)
		{
		bctx = .With(ctx, [blockArg: false])
		if ctx.blockArg is true
			return .FuncBody(node, bctx, curr, extra: .BlockParams(node, bctx, curr),
				lead: .Soft)
		extra = .BlockParams(node, bctx, curr)
		single = .Single?(node, ctx)
		sctx = .With(bctx, [constant: false, noSym: false, singleLine:])
		br = single ? .Line : .Hard
		sep = single ? .Semi : .Hard
		stmts = .StmtDocs(node, sctx, curr, allowBlank: not single)
		inner = Object()
		for (j = 0; j < stmts.Size(); ++j)
			{
			inner.Add(j is 0 ? br : sep)
			inner.Add(stmts[j])
			}
		inner.Add(br)
		inner.Add(.Cm.Leading(curr, node.end))
		inner.Add(.Text('}'))
		d = .Cat(.Text('{'), extra, .Nest(.Catl(inner)))
		return single ? .Group(d) : d
		}

	Param(node, ctx, curr)
		{
		parts = [.Text(node.name)]
		if node.hasdef
			{
			parts.Add(.Cm.Trailing(curr, node.pos + node.name.Size(),
				unusedParam: node.name is #unused))
			parts.Add(.Text(" = "))
			s = node.pos in (0, false) or node.end in (0, false)
				? false
				: .Style.ValueSrc(.Src[node.pos .. node.end].AfterFirst('='))
			parts.Add(.constant(node.defval, ctx, curr, :s))
			}
		return .Catl(parts)
		}
	}
