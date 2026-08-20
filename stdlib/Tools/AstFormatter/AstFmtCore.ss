// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
AstFmtDoc
	{
	With(ctx, mods)
		{
		c = ctx.Copy()
		for m in mods.Members()
			c[m] = mods[m]
		return c
		}

	Fmt(node, ctx, curr, allowBlankAfter = false)
		{
		if node is false
			return false
		unusedParam = node.type is #Param and node.name is #unused
		parts = Object()
		cctx = .LeadComments(node, ctx, curr, parts, :unusedParam)
		parts.Add(this[node.type](node, cctx, curr))
		.TrailComments(node, ctx, curr, parts, allowBlank: allowBlankAfter, :unusedParam)
		return .Catl(parts)
		}

	LeadComments(node, ctx, curr, parts, unusedParam = false)
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
				cctx = .With(ctx, [parentEnd: node.end])
			parts.Add(.Cm.Leading(curr, pos, :unusedParam))
			}
		return cctx
		}

	TrailComments(node, ctx, curr, parts, allowBlank = false, unusedParam = false)
		{
		// a scalar member is emitted as a verbatim slice; advance past it so any
		// in-slice comment is not re-emitted by Trailing
		if node.type is #Member and Type(node.value) isnt #AstNode and
			node.end not in (0, false)
			.Cm.SkipTo(curr, node.end)
		if false isnt node.pos and node.end not in (0, false)
			parts.Add(
				.Cm.Trailing(curr, node.end, ctx.parentEnd, :allowBlank, :unusedParam))
		}

	Params(params, ctx, curr, delims)
		{
		pctx = .With(ctx, [constant: false, noSym: false, singleLine:])
		docs = Object()
		for (i = 0; false isnt p = params[i]; ++i)
			docs.Add(.Fmt(p, pctx, curr))
		if docs.Empty?()
			return .Text(delims[0] $ delims[1])
		return .Group(.Cat(delims[0], .Nest(.Fillsep(docs, .Cat(',', .Line))), delims[1]))
		}

	BlockParams(node, ctx, curr)
		{
		params = node.params
		return params.size > 1 or (params.size is 1 and params[0].name isnt #it)
			? .Params(params, ctx, curr, "||")
			: false
		}

	FuncBody(node, ctx, curr, extra = false, pos2 = false, lead = false)
		{
		empty = .emptyFunc?(node, curr, pos2)
		single = empty or .Single?(node, ctx)
		bctx = .With(ctx, [constant: false, noSym: false, singleLine:])
		br = single ? .Line : .Hard
		sep = single ? .Semi : .Hard
		parts = [.bodyHead(single, lead), .Text('{')]
		if extra isnt false
			parts.Add(extra)
		.bodyOpen(parts, curr, empty, pos2)
		stmts = .StmtDocs(node, bctx, curr, allowBlank: not single)
		.bodyStmts(parts, stmts, br, sep)
		parts.Add(br)
		parts.Add(.Cm.LeadingToCode(curr))
		parts.Add(.Text('}'))
		d = .Nest(.Catl(parts))
		return single ? .Group(d) : d
		}

	// the break before '{': a multi-line body always breaks; a single-line
	// body hugs (.Soft, or the caller-supplied lead)
	bodyHead(single, lead)
		{
		if not single
			return .Hard
		return lead is false ? .Soft : lead
		}

	bodyOpen(parts, curr, empty, pos2)
		{
		if empty // Trailing would run past the close when { } are on one line
			.Cm.SkipTo(curr, pos2)
		else if pos2 isnt false
			parts.Add(.Cm.Trailing(curr, pos2))
		}

	bodyStmts(parts, stmts, br, sep)
		{
		for (j = 0; j < stmts.Size(); ++j)
			{
			parts.Add(j is 0 ? br : sep)
			parts.Add(stmts[j])
			}
		}

	// an empty body has nothing to lay out, so it hugs the signature as { }
	emptyFunc?(node, curr, pos2)
		{
		return node.type is #Function and node[0] is false and pos2 isnt false and
			.Cm.EmptyBraces?(curr, pos2)
		}

	Single?(node, ctx)
		{
		if ctx.singleLine isnt true or .BodyStyle is #expand
			return false
		return .BodyStyle is #fit
			? .Style.OkToSingleLine2(node)
			: .Style.OkToSingleLine(node)
		}

	StmtDocs(node, ctx, curr, allowBlank = true)
		{
		docs = Object()
		for (i = 0; false isnt stmt = node[i]; ++i)
			{
			last = false is node[i+1]
			d = .Fmt(stmt, .With(ctx, [lastStmt: last]), curr,
				allowBlankAfter: allowBlank and not last)
			if i+1 < node.size and node[i+1].type is #Compound and node[i+1].size is 0
				{
				d = .Cat(d, ";;", .Cm.Trailing(curr, node[i+1].pos))
				++i
				}
			docs.Add(d)
			}
		return docs
		}

	Rtrim(s)
		{
		while s.Size() > 0 and s[-1] in (' ', '\t', '\n', '\r')
			s = s[..-1]
		return s
		}
	}
