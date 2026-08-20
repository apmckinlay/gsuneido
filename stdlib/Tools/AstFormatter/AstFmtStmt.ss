// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// statement and control-flow shapes for AstFormatter
AstFmtExpr
	{
	Return(node, ctx, curr)
		{
		docs = Object()

		for (i = 0; false isnt e = node[i]; ++i)
			docs.Add(.Fmt(e, ctx, curr))

		kw = node.throw is true ? "return throw" : #return

		if docs.Empty?()
			return .Text(kw)

		return .Cat(kw $ ' ', .Seq(docs, .Text(", ")))
		}

	Throw(node, ctx, curr)
		{
		return .Cat("throw ", .Fmt(node.expr, ctx, curr))
		}

	TryCatch(node, ctx, curr)
		{
		parts = [.Text(#try)]
		parts.Add(.Cm.Trailing(curr, node.pos + 3))/*= "try".Size() */
		parts.Add(
			.ctlBody(node.try, ctx, curr, guard: node.catch is false ? false : #Catch))
		if false isnt c = node.catch
			{
			parts.Add(.Hard)
			parts.Add(.Text(#catch))
			if false isnt var = .Style.CatchVar(node)
				{
				v = .Fmt(var, ctx, curr)
				pat = node.catchpat is false ? "" : ", " $ Display(node.catchpat)
				parts.Add(.Cat(" (", v, pat, ')'))
				}
			parts.Add(.Cm.Trailing(curr, node.catchend))
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
			parts.Add(.Cm.Trailing(curr, node.elseend))
			f = node.f
			if false isnt inner = .Style.ElseChain(f, curr)
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
			.Cm.SkipTo(curr, end)
		return .Catl(parts)
		}

	cond(expr, ctx, curr)
		{
		return .Fmt(expr, .With(ctx, [cond:]), curr)
		}

	ctlBody(node, ctx, curr, guard = false, brace = false)
		{
		bctx = .With(ctx, [lastStmt: false, guard: false])
		if brace is false and false isnt stmt = .Style.Unbrace(node, guard)
			{
			d = .Fmt(stmt, bctx, curr)
			.Cm.SkipTo(curr, node.end)
			return .Nest(.Cat(.Hard, d))
			}
		// a bare loop/control body renders its own tail body in a separate
		// ctlBody; carry the else/catch guard down so that inner tail is not
		// unbraced into capturing our guard (dangling else/catch)
		if guard isnt false and node.type isnt #Compound
			bctx = .With(bctx, [:guard])
		return .Nest(.Cat(.Hard, .Fmt(node, bctx, curr)))
		}

	Switch(node, ctx, curr)
		{
		parts = [.Text(#switch)]
		e = node.expr is false ? false : .Style.SwitchExpr(node.expr)
		if e isnt false and (e.type isnt #Constant or e.value isnt true)
			{
			parts.Add(.Text(' '))
			parts.Add(.cond(e, ctx, curr))
			}
		parts.Add(.Cm.Trailing(curr, node.pos1))
		parts.Add(.Hard)
		parts.Add(.Text("\t{"))
		parts.Add(.Cm.Trailing(curr, node.pos2))
		for (i = 0; false isnt c = node[i]; ++i)
			{
			parts.Add(.Hard)
			parts.Add(.caseLabel(c, ctx, curr))
			parts.Add(.Cm.Trailing(curr, c.end))
			.caseBody(parts, c.body, ctx, curr)
			}
		if node.def isnt false
			{
			parts.Add(.Hard)
			parts.Add(.Text("default:"))
			parts.Add(.Cm.Trailing(curr, node.posdef))
			.caseBody(parts, node.def, ctx, curr)
			}
		parts.Add(.Hard)
		parts.Add(.Text("\t}"))
		return .Catl(parts)
		}

	caseLabel(c, ctx, curr)
		{
		vals = Object()
		for (j = 0; j < c.size; ++j)
			{
			if j > 0
				vals.Add(.Text(", "))
			vals.Add(.Fmt(c[j], ctx, curr))
			}
		return .Cat("case ", .Catl(vals), ':')
		}

	caseBody(parts, body, ctx, curr)
		{
		stmts = .StmtDocs(body, ctx, curr)
		if not stmts.Empty?()
			parts.Add(.Nest(.Cat(.Hard, .Seq(stmts, .Hard))))
		}

	Forever(node, ctx, curr)
		{
		return .Cat(#forever, .ctlBody(node.body, ctx, curr, guard: ctx.guard))
		}

	While(node, ctx, curr)
		{
		if node.cond.type is #Constant and node.cond.value is true
			return .Forever(node, ctx, curr)
		c = .cond(node.cond, ctx, curr)
		return .Cat("while ", c, .ctlBody(node.body, ctx, curr, guard: ctx.guard))
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
			header = .Rtrim(.Src[node.pos .. node.body.pos])
			.Cm.SkipTo(curr, node.body.pos)
			return .Cat(.Tok(header), .ctlBody(node.body, ctx, curr, guard: ctx.guard))
			}
		v = node.var
		if node.var2 isnt ""
			v $= ", " $ node.var2
		ex = .cond(node.expr, ctx, curr)
		return .Cat("for " $ v $ " in ", ex,
			.ctlBody(node.body, ctx, curr, guard: ctx.guard))
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
			parts.Add(.Expr(e, ctx, curr))
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
			parts.Add(.Expr(e, ctx, curr))
			sep = ", "
			}
		parts.Add(.Text(')'))
		return .Cat(.Catl(parts), .ctlBody(node.body, ctx, curr, guard: ctx.guard))
		}

	Compound(node, ctx, curr)
		{
		open = .Cm.Trailing(curr, node.pos + 1)
		stmts = .StmtDocs(node, ctx, curr)
		parts = [.Text('{'), open]
		if not stmts.Empty?()
			{
			parts.Add(.Hard)
			parts.Add(.Seq(stmts, .Hard))
			}
		parts.Add(.Hard)
		parts.Add(.Cm.Leading(curr, node.end - 1))
		parts.Add(.Text('}'))
		return .Catl(parts)
		}

	Break(unused, ctx/*unused*/, curr/*unused*/)
		{
		return .Text(#break)
		}

	Continue(unused, ctx/*unused*/, curr/*unused*/)
		{
		return .Text(#continue)
		}
	}
